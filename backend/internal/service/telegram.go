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
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramService публикует контент в канал/чат Telegram через Bot API.
type TelegramService interface {
	// PublishPost публикует пост в канал.
	// mediaPath — абсолютный путь к файлу на диске.
	// mediaType — "photo", "video" или "gif".
	// title, description, tags — текст подписи.
	PublishPost(mediaPath, mediaType, title, description, tags string) (int, error)
}

type telegramService struct {
	token   string
	chatID  string // ID канала (например -1001234567890) или @channel
	client  *http.Client
}

func NewTelegramService(botToken, chatID string) TelegramService {
	return &telegramService{
		token:  botToken,
		chatID: chatID,
		client: &http.Client{},
	}
}

func buildCaption(title, description, tags string) string {
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
	return strings.Join(parts, "\n\n")
}

// telegramResponse — общий ответ API с result.message_id.
type telegramResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
	Description string `json:"description,omitempty"`
}

func (s *telegramService) PublishPost(mediaPath, mediaType, title, description, tags string) (int, error) {
	if s.token == "" || s.chatID == "" {
		return 0, fmt.Errorf("telegram: bot token or chat_id not configured")
	}

	caption := buildCaption(title, description, tags)
	if len(caption) > 1024 {
		caption = caption[:1021] + "..."
	}

	var method string
	var fieldName string
	switch mediaType {
	case "photo":
		method = "sendPhoto"
		fieldName = "photo"
	case "video", "gif":
		method = "sendVideo"
		fieldName = "video"
	default:
		method = "sendDocument"
		fieldName = "document"
	}

	url := telegramAPIBase + s.token + "/" + method
	msgID, err := s.sendFile(url, fieldName, mediaPath, caption)
	if err != nil {
		return 0, err
	}
	return msgID, nil
}

func (s *telegramService) sendFile(apiURL, fieldName, localPath, caption string) (int, error) {
	f, err := os.Open(localPath)
	if err != nil {
		return 0, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// chat_id
	_ = w.WriteField("chat_id", s.chatID)
	if caption != "" {
		_ = w.WriteField("caption", caption)
	}

	// file
	name := filepath.Base(localPath)
	fw, err := w.CreateFormFile(fieldName, name)
	if err != nil {
		return 0, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(fw, f); err != nil {
		return 0, fmt.Errorf("copy file: %w", err)
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, apiURL, &buf)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("telegram request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response: %w", err)
	}

	var tgResp telegramResponse
	if err := json.Unmarshal(body, &tgResp); err != nil {
		return 0, fmt.Errorf("telegram decode: %w", err)
	}
	if !tgResp.OK {
		return 0, fmt.Errorf("telegram API error: %s", tgResp.Description)
	}
	return tgResp.Result.MessageID, nil
}
