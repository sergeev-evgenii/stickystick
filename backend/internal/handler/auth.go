package handler

import (
	"net/http"
	"sticky-stick/backend/internal/middleware"
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	analytics   service.AnalyticsService
}

func NewAuthHandler(authService service.AuthService, analytics service.AnalyticsService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		analytics:   analytics,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.analytics != nil {
		uid := user.ID
		_ = h.analytics.Log(middleware.ResolveClientIP(c), &uid, models.ActionRegister, "", 0, c.Request.UserAgent())
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if h.analytics != nil {
		uid := user.ID
		_ = h.analytics.Log(middleware.ResolveClientIP(c), &uid, models.ActionLogin, "", 0, c.Request.UserAgent())
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
