package models

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypeVideo MediaType = "video"
	MediaTypePhoto MediaType = "photo"
	MediaTypeGif   MediaType = "gif"
)

type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "pending"
	ModerationStatusApproved ModerationStatus = "approved"
	ModerationStatusRejected ModerationStatus = "rejected"
)

type Video struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	CategoryID  *uint          `json:"category_id" gorm:"index"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tags        string         `json:"tags" gorm:"type:text"` // JSON массив тегов или строка через запятую
	MediaType   MediaType      `json:"media_type" gorm:"type:varchar(20);default:'video'"`
	MediaURL          string           `json:"media_url" gorm:"not null"`                          // URL для видео, фото или гифки
	ThumbnailURL      string           `json:"thumbnail_url"`                                      // Превью для видео
	Duration          int              `json:"duration"`                                            // in seconds (для видео)
	Views             int              `json:"views" gorm:"default:0"`
	ModerationStatus  ModerationStatus `json:"moderation_status" gorm:"type:varchar(20);default:'pending';index"` // Статус модерации
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `json:"-" gorm:"index"`

	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Comments []Comment `json:"comments" gorm:"foreignKey:VideoID"`
	Likes    []Like    `json:"likes" gorm:"foreignKey:VideoID"`
}
