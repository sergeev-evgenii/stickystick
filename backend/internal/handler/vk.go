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

// VKHandler публикует контент из сервиса в группу ВКонтакте.
type VKHandler struct {
	vkService    service.VKService
	videoService service.VideoService
	mediaService service.MediaService
	userService  service.UserService
	settings     service.SettingsService
}

func NewVKHandler(vk service.VKService, video service.VideoService, media service.MediaService, user service.UserService, settings service.SettingsService) *VKHandler {
	return &VKHandler{
		vkService:    vk,
		videoService: video,
		mediaService: media,
		userService:  user,
		settings:     settings,
	}
}

// PublishVideoToVK godoc
// POST /api/v1/videos/:id/publish/vk
// Тело (опционально, JSON): { "comment": "доп. текст" }
// Требует авторизации администратора.
func (h *VKHandler) PublishVideoToVK(c *gin.Context) {
	// Проверяем, что запрос от администратора
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

	// Парсим ID видео
	idStr := c.Param("id")
	idParsed, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	// Загружаем видео из БД
	video, err := h.videoService.GetVideo(uint(idParsed), true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	// Опциональный дополнительный комментарий из тела запроса
	var body struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&body)

	// Формируем описание: описание видео + доп. комментарий
	description := video.Description
	comment := body.Comment
	if strings.TrimSpace(comment) != "" {
		if description != "" {
			description = description + "\n\n" + comment
		} else {
			description = comment
		}
	}
	// Для ВК: сначала заголовок поста (title), потом описание, потом ссылки (телега + проект).
	description = strings.TrimRight(description, "\n")
	footer := "мы в телеграм — " + telegramChannelURL + "\n" + "проект — " + projectURL
	if !strings.Contains(description, telegramChannelURL) && !strings.Contains(description, projectURL) {
		if strings.TrimSpace(description) != "" {
			description = description + "\n\n" + footer
		} else {
			description = footer
		}
	} else {
		// если какие-то ссылки уже есть — просто гарантируем, что обе присутствуют, без дублей
		description = ensureLinksFirst(description, []string{})
		if !strings.Contains(description, telegramChannelURL) {
			description = strings.TrimRight(description, "\n") + "\n" + "мы в телеграм — " + telegramChannelURL
		}
		if !strings.Contains(description, projectURL) {
			description = strings.TrimRight(description, "\n") + "\n" + "проект — " + projectURL
		}
	}

	// Определяем локальный путь к медиафайлу
	localPath, mediaType, err := h.resolveLocalPath(video.MediaURL, string(video.MediaType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("resolve media path: %v", err)})
		return
	}

	// Публикуем в ВК
	postID, err := h.vkService.PublishPost(localPath, mediaType, video.Title, description, video.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("vk publish: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id": postID,
		"message": "published to VK successfully",
	})
}

// resolveLocalPath преобразует MediaURL (вида /uploads/videos/xxx.mp4) в абсолютный путь на диске.
// Возвращает путь и нормализованный тип медиа ("photo", "video", "gif").
func (h *VKHandler) resolveLocalPath(mediaURL, mediaType string) (string, string, error) {
	// mediaService хранит uploadDir; получаем его через интерфейс
	localPath := h.mediaService.URLToPath(mediaURL)
	if localPath == "" {
		return "", "", fmt.Errorf("cannot resolve local path from url: %s", mediaURL)
	}

	// Нормализуем тип по расширению, если модель не задала
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
