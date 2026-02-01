package models

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	VideoID   uint           `json:"video_id" gorm:"not null;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Video Video `json:"-" gorm:"foreignKey:VideoID"`
	User  User  `json:"user" gorm:"foreignKey:UserID"`
}

// Unique constraint on VideoID and UserID
func (Like) TableName() string {
	return "likes"
}
