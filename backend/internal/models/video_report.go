package models

import (
	"time"

	"gorm.io/gorm"
)

// VideoReport — жалоба на видео (user_id пустой = анонимная жалоба)
type VideoReport struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	VideoID   uint           `json:"video_id" gorm:"not null;index"`
	UserID    *uint          `json:"user_id" gorm:"index"` // nil — анонимная жалоба
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
