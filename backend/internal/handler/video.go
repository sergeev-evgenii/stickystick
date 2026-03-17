package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/service"
	"sticky-stick/backend/internal/store"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService service.VideoService
	mediaService service.MediaService
	userService  service.UserService
	seenStore    *store.SeenStore
}

func NewVideoHandler(videoService service.VideoService, mediaService service.MediaService, userService service.UserService, seenStore *store.SeenStore) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
		mediaService: mediaService,
		userService:  userService,
		seenStore:    seenStore,
	}
}

// isAdmin проверяет, является ли пользователь администратором
func (h *VideoHandler) isAdmin(c *gin.Context) bool {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		return false
	}
	userID := userIDInterface.(uint)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		return false
	}
	return user.IsAdmin
}

func (h *VideoHandler) GetFeed(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	viewerKey := getViewerKey(c)
	seen := h.seenStore.GetSeen(viewerKey)
	excludeIDs := make([]uint, 0, len(seen))
	for id := range seen {
		excludeIDs = append(excludeIDs, id)
	}

	isAdmin := h.isAdmin(c)
	videos, err := h.videoService.GetFeed(limit, offset, isAdmin, excludeIDs, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Пользователь просмотрел всё — очищаем список и отдаём ленту в случайном порядке
	if len(videos) == 0 && len(excludeIDs) > 0 {
		h.seenStore.ClearSeen(viewerKey)
		videos, err = h.videoService.GetFeed(limit, offset, isAdmin, nil, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) GetVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	isAdmin := h.isAdmin(c)
	video, err := h.videoService.GetVideo(uint(id), isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	viewerKey := getViewerKey(c)
	h.seenStore.MarkSeen(viewerKey, []uint{video.ID})

	c.JSON(http.StatusOK, video)
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	var req struct {
		Title        string `json:"title" binding:"required"`
		Description  string `json:"description"`
		VideoURL     string `json:"video_url" binding:"required"` // Для обратной совместимости
		ThumbnailURL string `json:"thumbnail_url"`
		Duration     int    `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	// Используем новый метод UploadMedia для обратной совместимости
	video, err := h.videoService.UploadMedia(
		userID,
		req.Title,
		req.Description,
		"", // tags
		nil, // categoryID
		req.VideoURL,
		models.MediaTypeVideo,
		req.ThumbnailURL,
		req.Duration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, video)
}

// UploadMedia обрабатывает загрузку файла (фото, гифки или видео)
func (h *VideoHandler) UploadMedia(c *gin.Context) {
	// Получаем файл из формы
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file is required: %v", err)})
		return
	}

	// Проверяем размер файла (максимум 500MB)
	const maxFileSize = 500 << 20 // 500 MB
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file too large. Maximum size is %d MB", maxFileSize/(1<<20))})
		return
	}

	// Получаем остальные данные из формы
	title := c.PostForm("title")
	description := c.PostForm("description")
	tags := c.PostForm("tags")
	categoryIDStr := c.PostForm("category_id")
	
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	// Парсим categoryID
	var categoryID *uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			catID := uint(id)
			categoryID = &catID
		}
	}

	// Определяем тип медиа по расширению файла
	mediaType := h.mediaService.GetMediaType(file.Filename)
	
	// Валидация типа файла
	allowedTypes := []string{"video", "photo", "gif"}
	if err := h.mediaService.ValidateFileType(file, allowedTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем файл (SaveFile возвращает относительный путь videos/xxx.mp4 для БД)
	mediaURL, err := h.mediaService.SaveFile(file, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save file: %v", err)})
		return
	}

	// Если это видео, можно загрузить thumbnail отдельно
	var thumbnailURL string
	if mediaType == "video" {
		thumbnailFile, err := c.FormFile("thumbnail")
		if err == nil {
			thumbnailURL, err = h.mediaService.SaveFile(thumbnailFile, "photo")
			if err != nil {
				thumbnailURL = ""
			}
		}
	}

	// Получаем userID из контекста (установлен middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	// Создаем запись в БД
	video, err := h.videoService.UploadMedia(
		userID,
		title,
		description,
		tags,
		categoryID,
		mediaURL,
		models.MediaType(mediaType),
		thumbnailURL,
		0, // duration можно определить позже или из метаданных
	)
	if err != nil {
		// Если не удалось создать запись, удаляем загруженный файл
		_ = h.mediaService.DeleteFile(mediaURL)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, video)
}

func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	// Получаем userID из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	if err := h.videoService.DeleteVideo(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
}

func (h *VideoHandler) LikeVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	// Получаем userID из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	if err := h.videoService.LikeVideo(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video liked"})
}

func (h *VideoHandler) UnlikeVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	// Получаем userID из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	if err := h.videoService.UnlikeVideo(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video unliked"})
}

func (h *VideoHandler) AddComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	comment, err := h.videoService.AddComment(uint(id), userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetPendingModeration возвращает список видео на модерации (только для админов)
func (h *VideoHandler) GetPendingModeration(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	videos, err := h.videoService.GetPendingModeration(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}

// ModerateVideo одобряет или отклоняет видео (только для админов)
func (h *VideoHandler) ModerateVideo(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"` // "approved" или "rejected"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status models.ModerationStatus
	switch req.Status {
	case "approved":
		status = models.ModerationStatusApproved
	case "rejected":
		status = models.ModerationStatusRejected
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. Use 'approved' or 'rejected'"})
		return
	}

	if err := h.videoService.ModerateVideo(uint(id), status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video moderated successfully", "status": status})
}

func (h *VideoHandler) GetApproved(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	limit, offset := 50, 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}
	videos, err := h.videoService.GetApproved(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) GetHidden(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	limit, offset := 50, 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}
	videos, err := h.videoService.GetHidden(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) HideVideo(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}
	if err := h.videoService.HideVideo(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "hidden"})
}

func (h *VideoHandler) UnhideVideo(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}
	if err := h.videoService.UnhideVideo(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "visible"})
}

func (h *VideoHandler) UpdateVideoFields(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.videoService.UpdateVideoFields(uint(id), req.Title, req.Description, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
