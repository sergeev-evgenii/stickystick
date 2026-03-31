package main

import (
	"log"
	"net/http"

	"sticky-stick/backend/internal/config"
	"sticky-stick/backend/internal/handler"
	"sticky-stick/backend/internal/middleware"
	"sticky-stick/backend/internal/repository"
	"sticky-stick/backend/internal/service"
	"sticky-stick/backend/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos, cfg)

	// Хранилище «просмотренных» видео по зрителю (пока в памяти)
	seenStore := store.NewSeenStore()

	// Initialize handlers
	handlers := handler.NewHandlers(services, seenStore)

	// Setup router
	router := setupRouter(handlers, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(h *handler.Handlers, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Игнорируем запросы от Chrome DevTools
	router.GET("/.well-known/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
		c.Abort()
	})

	// CORS middleware - должен быть первым и применяться ко всем роутам
	router.Use(corsMiddleware())
	// IP для аналитики и логов (X-Forwarded-For / RemoteAddr)
	router.Use(middleware.ClientIPMiddleware())

	// Увеличиваем лимит размера загружаемого файла до 500MB
	router.MaxMultipartMemory = 500 << 20 // 500 MB

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve static files (uploads)
	router.Static("/uploads", cfg.UploadDir)

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id", h.User.GetProfile)
			users.PUT("/:id", middleware.AuthMiddleware(cfg), h.User.UpdateProfile) // Требует авторизации
		}

		// Video routes
		videos := api.Group("/videos")
		{
			videos.GET("", middleware.OptionalAuthMiddleware(cfg), h.Video.GetFeed)
			videos.POST("", middleware.AuthMiddleware(cfg), h.Video.UploadVideo)
			videos.POST("/upload", middleware.AuthMiddleware(cfg), h.Video.UploadMedia)

			// Модерация — статические пути ОБЯЗАТЕЛЬНО до /:id
			videos.GET("/moderation/pending", middleware.AuthMiddleware(cfg), h.Video.GetPendingModeration)
			videos.GET("/moderation/approved", middleware.AuthMiddleware(cfg), h.Video.GetApproved)
			videos.GET("/moderation/hidden", middleware.AuthMiddleware(cfg), h.Video.GetHidden)

			// Параметрические пути
			videos.GET("/:id", middleware.OptionalAuthMiddleware(cfg), h.Video.GetVideo)
			videos.DELETE("/:id", middleware.AuthMiddleware(cfg), h.Video.DeleteVideo)
			videos.PUT("/:id", middleware.AuthMiddleware(cfg), h.Video.UpdateVideoFields)
			videos.POST("/:id/like", middleware.AuthMiddleware(cfg), h.Video.LikeVideo)
			videos.DELETE("/:id/like", middleware.AuthMiddleware(cfg), h.Video.UnlikeVideo)
			videos.POST("/:id/comment", middleware.AuthMiddleware(cfg), h.Video.AddComment)
			videos.POST("/:id/moderate", middleware.AuthMiddleware(cfg), h.Video.ModerateVideo)
			videos.POST("/:id/hide", middleware.AuthMiddleware(cfg), h.Video.HideVideo)
			videos.POST("/:id/unhide", middleware.AuthMiddleware(cfg), h.Video.UnhideVideo)
			videos.POST("/:id/publish/vk", middleware.AuthMiddleware(cfg), h.VK.PublishVideoToVK)
			videos.POST("/:id/publish/telegram", middleware.AuthMiddleware(cfg), h.Telegram.PublishVideoToTelegram)
			videos.POST("/:id/publish/max", middleware.AuthMiddleware(cfg), h.Max.PublishVideoToMax)
		}

		// Admin routes
		admin := api.Group("/admin")
		{
			admin.GET("/analytics", middleware.AuthMiddleware(cfg), h.Admin.GetAnalytics)
		}

		// Публичная аналитика: логирование нажатия «Сгенерировать своё видео» (IP, кол-во переходов)
		api.POST("/analytics/generate-video-click", middleware.OptionalAuthMiddleware(cfg), h.Admin.LogGenerateVideoClick)

		// Settings — GET публичный, PATCH только админ
		api.GET("/settings", h.Settings.GetPublic)
		api.PATCH("/settings", middleware.AuthMiddleware(cfg), h.Settings.UpdateShowViewCount)

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", h.Category.GetAll)                    // Публичный доступ
			categories.GET("/:id", h.Category.GetByID)              // Публичный доступ
			categories.POST("", middleware.AuthMiddleware(cfg), h.Category.Create)    // Требует авторизации (админ)
			categories.PUT("/:id", middleware.AuthMiddleware(cfg), h.Category.Update) // Требует авторизации (админ)
			categories.DELETE("/:id", middleware.AuthMiddleware(cfg), h.Category.Delete) // Требует авторизации (админ)
		}
	}

	return router
}

// corsMiddleware обрабатывает CORS заголовки
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Разрешаем все origins для разработки
		// Если origin указан, используем его (обязательно для credentials)
		// Если нет origin, используем * (но без credentials)
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Обрабатываем preflight запрос (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}
