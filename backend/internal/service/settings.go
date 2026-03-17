package service

import (
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/repository"
)

type SettingsService interface {
	GetPublic() (*models.SiteSettings, error)
	SetShowViewCount(value bool) error
	SetDefaultPublishTexts(vk, telegram, max *string) error
}

type settingsService struct {
	repo repository.SettingsRepository
}

func NewSettingsService(repo repository.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

func (s *settingsService) GetPublic() (*models.SiteSettings, error) {
	return s.repo.Get()
}

func (s *settingsService) SetShowViewCount(value bool) error {
	return s.repo.UpdateShowViewCount(value)
}

func (s *settingsService) SetDefaultPublishTexts(vk, telegram, max *string) error {
	return s.repo.UpdateDefaults(vk, telegram, max)
}
