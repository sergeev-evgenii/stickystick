package main

import (
	"log"
	"net/http"
	"os"

	"sticky-stick/backend/internal/config"
	"sticky-stick/backend/internal/handler"
	"sticky-stick/backend/internal/middleware"
	"sticky-stick/backend/internal/repository"
	"sticky-stick/backend/internal/service"

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

	// Initialize handlers
	handlers := handler.NewHandlers(services)

	// Setup router
	router := setupRouter(handlers, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
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
			videos.GET("", middleware.OptionalAuthMiddleware(cfg), h.Video.GetFeed)                    // Публичный доступ (с опциональной авторизацией для админов)
			videos.GET("/:id", middleware.OptionalAuthMiddleware(cfg), h.Video.GetVideo)                // Публичный доступ (с опциональной авторизацией для админов)
			videos.POST("", middleware.AuthMiddleware(cfg), h.Video.UploadVideo)           // Требует авторизации
			videos.POST("/upload", middleware.AuthMiddleware(cfg), h.Video.UploadMedia)    // Требует авторизации
			videos.DELETE("/:id", middleware.AuthMiddleware(cfg), h.Video.DeleteVideo)    // Требует авторизации
			videos.POST("/:id/like", middleware.AuthMiddleware(cfg), h.Video.LikeVideo)    // Требует авторизации
			videos.DELETE("/:id/like", middleware.AuthMiddleware(cfg), h.Video.UnlikeVideo) // Требует авторизации
			videos.POST("/:id/comment", middleware.AuthMiddleware(cfg), h.Video.AddComment) // Требует авторизации
			
			// Модерация (только для админов)
			videos.GET("/moderation/pending", middleware.AuthMiddleware(cfg), h.Video.GetPendingModeration)
			videos.POST("/:id/moderate", middleware.AuthMiddleware(cfg), h.Video.ModerateVideo)
		}

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
