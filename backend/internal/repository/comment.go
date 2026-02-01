package repository

import (
	"sticky-stick/backend/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByVideoID(videoID uint, limit, offset int) ([]models.Comment, error)
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByVideoID(videoID uint, limit, offset int) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("video_id = ?", videoID).Preload("User").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&comments).Error
	return comments, err
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
