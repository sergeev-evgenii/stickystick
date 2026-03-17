package handler

import (
	"net/http"
	"strconv"
	"time"

	"sticky-stick/backend/internal/middleware"
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	analyticsService service.AnalyticsService
	userService      service.UserService
}

func NewAdminHandler(analyticsService service.AnalyticsService, userService service.UserService) *AdminHandler {
	return &AdminHandler{
		analyticsService: analyticsService,
		userService:      userService,
	}
}

func (h *AdminHandler) isAdmin(c *gin.Context) bool {
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

// GetAnalytics возвращает статистику и последнюю активность (только для админов).
// Query: since=24h | 7d | 30d (по умолчанию 24h), limit=100, offset=0
func (h *AdminHandler) GetAnalytics(c *gin.Context) {
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	sinceStr := c.DefaultQuery("since", "24h")
	var since time.Time
	switch sinceStr {
	case "7d":
		since = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		since = time.Now().Add(-30 * 24 * time.Hour)
	default:
		since = time.Now().Add(-24 * time.Hour)
	}

	limit := 100
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 500 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	stats, err := h.analyticsService.GetStats(since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activity, err := h.analyticsService.GetRecentActivity(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обогащаем активность пользователями из основной БД
	var userIDs []uint
	for i := range activity {
		if activity[i].UserID != nil {
			userIDs = append(userIDs, *activity[i].UserID)
		}
	}
	userMap := make(map[uint]*models.User)
	if len(userIDs) > 0 {
		users, _ := h.userService.GetByIDs(userIDs)
		for _, u := range users {
			userMap[u.ID] = u
		}
	}
	for i := range activity {
		if activity[i].UserID != nil {
			activity[i].User = userMap[*activity[i].UserID]
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"since":    since.Format(time.RFC3339),
		"period":   sinceStr,
		"stats":    stats,
		"activity": activity,
	})
}

// LogGenerateVideoClick — публичный endpoint: логирует нажатие «Сгенерировать своё видео» (IP, user_id если есть) в аналитику.
// Требует ClientIPMiddleware и опционально OptionalAuthMiddleware на роуте.
func (h *AdminHandler) LogGenerateVideoClick(c *gin.Context) {
	ip := middleware.GetClientIP(c)
	userAgent := c.Request.UserAgent()
	var userID *uint
	if idVal, exists := c.Get("userID"); exists {
		if id, ok := idVal.(uint); ok {
			userID = &id
		}
	}
	if err := h.analyticsService.Log(ip, userID, models.ActionGenerateVideoClick, "", 0, userAgent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
