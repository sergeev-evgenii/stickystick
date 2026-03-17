package models

// SiteSettings — одна запись с настройками сайта (id=1).
type SiteSettings struct {
	ID uint `json:"id" gorm:"primaryKey"`

	ShowViewCount bool `json:"show_view_count" gorm:"default:true;not null"`

	DefaultPublishVK       string `json:"default_publish_vk" gorm:"type:text;default:'';not null"`
	DefaultPublishTelegram string `json:"default_publish_telegram" gorm:"type:text;default:'';not null"`
	DefaultPublishMax      string `json:"default_publish_max" gorm:"type:text;default:'';not null"`
}

// TableName задаёт имя таблицы.
func (SiteSettings) TableName() string {
	return "site_settings"
}
