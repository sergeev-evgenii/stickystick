package repository

import (
	"sticky-stick/backend/internal/models"

	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(like *models.Like) error
	Delete(videoID, userID uint) error
	Exists(videoID, userID uint) (bool, error)
	CountByVideoID(videoID uint) (int64, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(videoID, userID uint) error {
	return r.db.Where("video_id = ? AND user_id = ?", videoID, userID).
		Delete(&models.Like{}).Error
}

func (r *likeRepository) Exists(videoID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).
		Where("video_id = ? AND user_id = ?", videoID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *likeRepository) CountByVideoID(videoID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}
