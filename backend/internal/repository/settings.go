package repository

import (
	"sticky-stick/backend/internal/models"

	"gorm.io/gorm"
)

type SettingsRepository interface {
	Get() (*models.SiteSettings, error)
	UpdateShowViewCount(value bool) error
	UpdateDefaults(vk, telegram, max *string) error
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) Get() (*models.SiteSettings, error) {
	var s models.SiteSettings
	err := r.db.Where("id = ?", 1).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *settingsRepository) UpdateShowViewCount(value bool) error {
	return r.db.Model(&models.SiteSettings{}).Where("id = ?", 1).Update("show_view_count", value).Error
}

func (r *settingsRepository) UpdateDefaults(vk, telegram, max *string) error {
	updates := map[string]interface{}{}
	if vk != nil {
		updates["default_publish_vk"] = *vk
	}
	if telegram != nil {
		updates["default_publish_telegram"] = *telegram
	}
	if max != nil {
		updates["default_publish_max"] = *max
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&models.SiteSettings{}).Where("id = ?", 1).Updates(updates).Error
}
