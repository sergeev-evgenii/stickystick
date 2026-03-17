package handler

import (
	"sticky-stick/backend/internal/service"
	"sticky-stick/backend/internal/store"
)

type Handlers struct {
	Auth      *AuthHandler
	User      *UserHandler
	Video     *VideoHandler
	Category  *CategoryHandler
	Admin     *AdminHandler
	VK        *VKHandler
	Telegram  *TelegramHandler
	Max       *MaxHandler
	Settings  *SettingsHandler
}

func NewHandlers(services *service.Services, seenStore *store.SeenStore) *Handlers {
	return &Handlers{
		Auth:     NewAuthHandler(services.Auth),
		User:     NewUserHandler(services.User),
		Video:    NewVideoHandler(services.Video, services.Media, services.User, seenStore),
		Category: NewCategoryHandler(services.Category),
		Admin:    NewAdminHandler(services.Analytics, services.User),
		VK:       NewVKHandler(services.VK, services.Video, services.Media, services.User, services.Settings),
		Telegram: NewTelegramHandler(services.Telegram, services.Video, services.Media, services.User, services.Settings),
		Max:      NewMaxHandler(services.Max, services.Video, services.Media, services.User, services.Settings),
		Settings: NewSettingsHandler(services.Settings, services.User),
	}
}
