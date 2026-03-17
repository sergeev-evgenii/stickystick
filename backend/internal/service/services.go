package service

import (
	"sticky-stick/backend/internal/config"
	"sticky-stick/backend/internal/repository"
)

type Services struct {
	Auth      AuthService
	User      UserService
	Video     VideoService
	Media     MediaService
	Category  CategoryService
	Analytics AnalyticsService
	VK        VKService
	Telegram  TelegramService
	Max       MaxService
	Settings  SettingsService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:      NewAuthService(repos.User, cfg),
		User:      NewUserService(repos.User),
		Video:     NewVideoService(repos.Video, repos.Comment, repos.Like),
		Media:     NewMediaService(cfg.UploadDir, cfg.BaseURL),
		Category:  NewCategoryService(repos.Category),
		Analytics: NewAnalyticsService(repos.ActivityLog, repos.Video),
		VK:        NewVKService(cfg.VKToken, cfg.VKGroupID),
		Telegram:  NewTelegramService(cfg.TelegramBotToken, cfg.TelegramChatID),
		Max:       NewMaxService(cfg.MaxBotToken, cfg.MaxChatID),
		Settings:  NewSettingsService(repos.Settings),
	}
}
