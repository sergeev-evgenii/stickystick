package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Avatar    string         `json:"avatar"`
	Bio       string         `json:"bio"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"` // Флаг администратора
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Videos   []Video   `json:"-" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"-" gorm:"foreignKey:UserID"`
	Likes    []Like    `json:"-" gorm:"foreignKey:UserID"`
}
