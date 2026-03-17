package models

import (
	"time"

	"gorm.io/gorm"
)

// Action types for activity log
const (
	ActionLogin     = "login"
	ActionRegister  = "register"
	ActionVideoView = "video_view"
	ActionLike      = "like"
	ActionUnlike    = "unlike"
	ActionUpload    = "upload"
	ActionFeedView  = "feed_view"
	ActionGenerateVideoClick = "generate_video_click"
)

type ActivityLog struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	IP           string         `json:"ip" gorm:"not null;index"`
	UserID       *uint          `json:"user_id" gorm:"index"`
	Action       string         `json:"action" gorm:"not null;index"`
	ResourceType string         `json:"resource_type"`
	ResourceID   uint           `json:"resource_id"`
	UserAgent    string         `json:"user_agent" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// User заполняется при отдаче из основной БД (не связь в БД аналитики — users в другой БД)
	User *User `json:"user,omitempty" gorm:"-"`
}
