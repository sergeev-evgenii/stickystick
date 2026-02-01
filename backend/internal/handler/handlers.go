package handler

import "sticky-stick/backend/internal/service"

type Handlers struct {
	Auth     *AuthHandler
	User     *UserHandler
	Video    *VideoHandler
	Category *CategoryHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:     NewAuthHandler(services.Auth),
		User:     NewUserHandler(services.User),
		Video:    NewVideoHandler(services.Video, services.Media, services.User),
		Category: NewCategoryHandler(services.Category),
	}
}
