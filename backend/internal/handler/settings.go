package handler

import (
	"net/http"

	"sticky-stick/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsService service.SettingsService
	userService     service.UserService
}

func NewSettingsHandler(settingsService service.SettingsService, userService service.UserService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		userService:     userService,
	}
}

// GetPublic возвращает публичные настройки (без авторизации).
func (h *SettingsHandler) GetPublic(c *gin.Context) {
	s, err := h.settingsService.GetPublic()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"show_view_count":            s.ShowViewCount,
		"default_publish_vk":         s.DefaultPublishVK,
		"default_publish_telegram":   s.DefaultPublishTelegram,
		"default_publish_max":        s.DefaultPublishMax,
	})
}

// UpdateShowViewCount обновляет настройки (только админ).
func (h *SettingsHandler) UpdateShowViewCount(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)
	user, err := h.userService.GetProfile(userID)
	if err != nil || !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var body struct {
		ShowViewCount          *bool   `json:"show_view_count"`
		DefaultPublishVK       *string `json:"default_publish_vk"`
		DefaultPublishTelegram *string `json:"default_publish_telegram"`
		DefaultPublishMax      *string `json:"default_publish_max"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if body.ShowViewCount != nil {
		if err := h.settingsService.SetShowViewCount(*body.ShowViewCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if err := h.settingsService.SetDefaultPublishTexts(body.DefaultPublishVK, body.DefaultPublishTelegram, body.DefaultPublishMax); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s, err := h.settingsService.GetPublic()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"show_view_count":            s.ShowViewCount,
		"default_publish_vk":         s.DefaultPublishVK,
		"default_publish_telegram":   s.DefaultPublishTelegram,
		"default_publish_max":        s.DefaultPublishMax,
	})
}
