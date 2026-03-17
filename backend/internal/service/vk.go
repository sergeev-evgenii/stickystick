package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const vkAPIVersion = "5.199"
const vkAPIBase = "https://api.vk.com/method"

// VKService публикует контент в группу ВКонтакте.
type VKService interface {
	// PublishPost публикует пост в группу.
	// mediaPath — абсолютный путь к файлу на диске (может быть пустым).
	// mediaType — "photo", "video" или "" (только текст).
	// title, description, tags — текстовые поля.
	PublishPost(mediaPath, mediaType, title, description, tags string) (int, error)
}

type vkService struct {
	token   string
	groupID string // без минуса, только цифры
	client  *http.Client
}

func NewVKService(token, groupID string) VKService {
	return &vkService{
		token:   token,
		groupID: groupID,
		client:  &http.Client{},
	}
}

// ownerID возвращает groupID со знаком минус (для wall.post owner_id = -groupID).
func (s *vkService) ownerID() string {
	return "-" + s.groupID
}

// buildMessage собирает текст поста из полей видео.
func buildMessage(title, description, tags string) string {
	var parts []string
	if title != "" {
		parts = append(parts, title)
	}
	if description != "" {
		parts = append(parts, description)
	}
	if tags != "" {
		// теги могут быть через запятую — добавляем # перед каждым
		rawTags := strings.Split(tags, ",")
		var hashTags []string
		for _, t := range rawTags {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			if !strings.HasPrefix(t, "#") {
				t = "#" + t
			}
			hashTags = append(hashTags, t)
		}
		if len(hashTags) > 0 {
			parts = append(parts, strings.Join(hashTags, " "))
		}
	}
	return strings.Join(parts, "\n\n")
}

// vkError описывает ошибку из VK API.
type vkError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e vkError) Error() string {
	return fmt.Sprintf("VK API error %d: %s", e.Code, e.Message)
}

// vkCall выполняет GET-запрос к VK API и декодирует ответ в out.
func (s *vkService) vkCall(method string, params url.Values, out interface{}) error {
	params.Set("access_token", s.token)
	params.Set("v", vkAPIVersion)

	resp, err := s.client.Get(vkAPIBase + "/" + method + "?" + params.Encode())
	if err != nil {
		return fmt.Errorf("vk http get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("vk read body: %w", err)
	}

	// Проверяем наличие поля error в ответе
	var errWrap struct {
		Error *vkError `json:"error"`
	}
	if err := json.Unmarshal(body, &errWrap); err == nil && errWrap.Error != nil {
		return errWrap.Error
	}

	// Декодируем в целевую структуру
	var wrapper struct {
		Response json.RawMessage `json:"response"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return fmt.Errorf("vk decode wrapper: %w", err)
	}
	if out != nil {
		if err := json.Unmarshal(wrapper.Response, out); err != nil {
			return fmt.Errorf("vk decode response: %w", err)
		}
	}
	return nil
}

// PublishPost — основная точка входа.
func (s *vkService) PublishPost(mediaPath, mediaType, title, description, tags string) (int, error) {
	message := buildMessage(title, description, tags)

	var attachments string
	var err error

	switch mediaType {
	case "photo":
		attachments, err = s.uploadPhoto(mediaPath)
		if err != nil {
			return 0, fmt.Errorf("upload photo: %w", err)
		}
	case "video", "gif":
		attachments, err = s.uploadVideo(mediaPath, title, description)
		if err != nil {
			return 0, fmt.Errorf("upload video: %w", err)
		}
	}

	return s.wallPost(message, attachments)
}

// --- Photo upload -----------------------------------------------------------

type photoUploadServer struct {
	UploadURL string `json:"upload_url"`
	AlbumID   int    `json:"album_id"`
	UserID    int    `json:"user_id"`
}

type photoUploadResponse struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type savedPhoto struct {
	ID      int `json:"id"`
	OwnerID int `json:"owner_id"`
}

func (s *vkService) uploadPhoto(localPath string) (string, error) {
	// 1. Получаем URL для загрузки
	var server photoUploadServer
	if err := s.vkCall("photos.getWallUploadServer", url.Values{
		"group_id": {s.groupID},
	}, &server); err != nil {
		return "", err
	}

	// 2. Загружаем файл на сервер VK
	uploadResp, err := s.multipartUpload(server.UploadURL, "photo", localPath)
	if err != nil {
		return "", fmt.Errorf("multipart upload photo: %w", err)
	}

	var up photoUploadResponse
	if err := json.Unmarshal(uploadResp, &up); err != nil {
		return "", fmt.Errorf("decode photo upload response: %w", err)
	}

	// 3. Сохраняем фото
	var saved []savedPhoto
	if err := s.vkCall("photos.saveWallPhoto", url.Values{
		"group_id": {s.groupID},
		"server":   {fmt.Sprintf("%d", up.Server)},
		"photo":    {up.Photo},
		"hash":     {up.Hash},
	}, &saved); err != nil {
		return "", err
	}
	if len(saved) == 0 {
		return "", fmt.Errorf("photos.saveWallPhoto returned empty list")
	}

	return fmt.Sprintf("photo%d_%d", saved[0].OwnerID, saved[0].ID), nil
}

// --- Video upload -----------------------------------------------------------

type videoSaveResponse struct {
	UploadURL string `json:"upload_url"`
	VideoID   int    `json:"video_id"`
	OwnerID   int    `json:"owner_id"`
}

func (s *vkService) uploadVideo(localPath, title, description string) (string, error) {
	// 1. Запрашиваем URL для загрузки видео
	params := url.Values{
		"group_id":    {s.groupID},
		"name":        {title},
		"description": {description},
		"wallpost":    {"0"}, // сами сделаем wall.post
	}
	var saveResp videoSaveResponse
	if err := s.vkCall("video.save", params, &saveResp); err != nil {
		return "", err
	}

	// 2. Загружаем файл на upload_url (PUT/POST с файлом)
	if _, err := s.multipartUpload(saveResp.UploadURL, "video_file", localPath); err != nil {
		return "", fmt.Errorf("multipart upload video: %w", err)
	}

	return fmt.Sprintf("video%d_%d", saveResp.OwnerID, saveResp.VideoID), nil
}

// --- wall.post --------------------------------------------------------------

type wallPostResponse struct {
	PostID int `json:"post_id"`
}

func (s *vkService) wallPost(message, attachments string) (int, error) {
	params := url.Values{
		"owner_id":  {s.ownerID()},
		"from_group": {"1"},
		"message":   {message},
	}
	if attachments != "" {
		params.Set("attachments", attachments)
	}

	var resp wallPostResponse
	if err := s.vkCall("wall.post", params, &resp); err != nil {
		return 0, err
	}
	return resp.PostID, nil
}

// --- helpers ----------------------------------------------------------------

// multipartUpload загружает файл multipart/form-data на uploadURL и возвращает тело ответа.
func (s *vkService) multipartUpload(uploadURL, fieldName, localPath string) (json.RawMessage, error) {
	f, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile(fieldName, filepath.Base(localPath))
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(fw, f); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, uploadURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http post upload: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read upload response: %w", err)
	}
	return body, nil
}
