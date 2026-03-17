package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxAPIBase = "https://platform-api.max.ru"

// MaxService публикует контент в чат/канал мессенджера Max через Platform API.
type MaxService interface {
	// PublishPost загружает файл в Max и отправляет сообщение с вложением.
	// mediaPath — абсолютный путь к файлу на диске.
	// mediaType — "photo", "video" или "gif".
	PublishPost(mediaPath, mediaType, title, description, tags string) (string, error)
}

type maxService struct {
	token  string
	chatID string
	client *http.Client
}

func NewMaxService(botToken, chatID string) MaxService {
	return &maxService{
		token:  botToken,
		chatID: chatID,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func buildMaxCaption(title, description, tags string) string {
	var parts []string
	if title != "" {
		parts = append(parts, title)
	}
	if description != "" {
		parts = append(parts, description)
	}
	if tags != "" {
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
	caption := strings.Join(parts, "\n\n")
	if len(caption) > 4000 {
		caption = caption[:3997] + "..."
	}
	return caption
}

// getUploadURL запрашивает URL для загрузки файла (type=image|video|file).
func (s *maxService) getUploadURL(mediaType string) (string, error) {
	maxType := "video"
	switch mediaType {
	case "photo", "gif":
		maxType = "image"
	case "video":
		maxType = "video"
	default:
		maxType = "file"
	}
	req, err := http.NewRequest(http.MethodPost, maxAPIBase+"/uploads?type="+maxType, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.token)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("max get upload url: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("max upload url: %s %s", resp.Status, string(body))
	}
	var out struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("max parse upload url: %w", err)
	}
	if out.URL == "" {
		return "", fmt.Errorf("max: empty upload url in response")
	}
	return out.URL, nil
}

// uploadFile загружает файл по URL и возвращает token вложения.
func (s *maxService) uploadFile(uploadURL, localPath string) (string, error) {
	f, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	name := filepath.Base(localPath)
	fw, err := w.CreateFormFile("data", name)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return "", err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, uploadURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", s.token)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("max upload file: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("max upload file: %s %s", resp.Status, string(body))
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("max parse upload response: %w", err)
	}
	if out.Token == "" {
		return "", fmt.Errorf("max: empty token in upload response")
	}
	return out.Token, nil
}

// sendMessage отправляет сообщение с вложением в чат.
func (s *maxService) sendMessage(caption, attachmentType, token string) error {
	payload := map[string]interface{}{
		"text": caption,
		"attachments": []map[string]interface{}{
			{
				"type": attachmentType,
				"payload": map[string]string{"token": token},
			},
		},
	}
	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, maxAPIBase+"/messages?chat_id="+s.chatID, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("max send message: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("max send message: %s %s", resp.Status, string(respBody))
	}
	return nil
}

func (s *maxService) PublishPost(mediaPath, mediaType, title, description, tags string) (string, error) {
	if s.token == "" || s.chatID == "" {
		return "", fmt.Errorf("max: bot token or chat_id not configured")
	}
	caption := buildMaxCaption(title, description, tags)

	maxUploadType := "video"
	maxAttachmentType := "video"
	switch mediaType {
	case "photo", "gif":
		maxUploadType = "image"
		maxAttachmentType = "image"
	case "video":
		maxUploadType = "video"
		maxAttachmentType = "video"
	default:
		maxUploadType = "file"
		maxAttachmentType = "file"
	}

	uploadURL, err := s.getUploadURL(maxUploadType)
	if err != nil {
		return "", err
	}
	token, err := s.uploadFile(uploadURL, mediaPath)
	if err != nil {
		return "", err
	}
	// MAX может вернуть attachment.not.ready при слишком быстрой отправке — небольшая пауза
	time.Sleep(1 * time.Second)
	if err := s.sendMessage(caption, maxAttachmentType, token); err != nil {
		return "", err
	}
	return token, nil
}
