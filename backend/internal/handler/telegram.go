package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"sticky-stick/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// TelegramHandler публикует контент в канал Telegram.
type TelegramHandler struct {
	telegramService service.TelegramService
	videoService    service.VideoService
	mediaService    service.MediaService
	userService    service.UserService
	settings       service.SettingsService
}

func NewTelegramHandler(telegram service.TelegramService, video service.VideoService, media service.MediaService, user service.UserService, settings service.SettingsService) *TelegramHandler {
	return &TelegramHandler{
		telegramService: telegram,
		videoService:    video,
		mediaService:    media,
		userService:     user,
		settings:        settings,
	}
}

// PublishVideoToTelegram godoc
// POST /api/v1/videos/:id/publish/telegram
// Тело (опционально, JSON): { "comment": "доп. текст" }
// Требует авторизации администратора.
func (h *TelegramHandler) PublishVideoToTelegram(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)
	user, err := h.userService.GetProfile(userID)
	if err != nil || !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	idStr := c.Param("id")
	idParsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	video, err := h.videoService.GetVideo(uint(idParsed), true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	var body struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&body)

	description := video.Description
	comment := body.Comment
	if strings.TrimSpace(comment) == "" && h.settings != nil {
		if s, err := h.settings.GetPublic(); err == nil {
			comment = s.DefaultPublishTelegram
		}
	}
	if strings.TrimSpace(comment) != "" {
		if description != "" {
			description = description + "\n\n" + comment
		} else {
			description = comment
		}
	}

	localPath, mediaType, err := h.resolveLocalPath(video.MediaURL, string(video.MediaType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("resolve media path: %v", err)})
		return
	}

	messageID, err := h.telegramService.PublishPost(localPath, mediaType, video.Title, description, video.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("telegram publish: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": messageID,
		"message":    "published to Telegram successfully",
	})
}

func (h *TelegramHandler) resolveLocalPath(mediaURL, mediaType string) (string, string, error) {
	localPath := h.mediaService.URLToPath(mediaURL)
	if localPath == "" {
		return "", "", fmt.Errorf("cannot resolve local path from url: %s", mediaURL)
	}
	if mediaType == "" {
		ext := strings.ToLower(filepath.Ext(localPath))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
			mediaType = "photo"
		case ".gif":
			mediaType = "gif"
		default:
			mediaType = "video"
		}
	}
	return localPath, mediaType, nil
}
