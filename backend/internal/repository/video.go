package repository

import (
	"sticky-stick/backend/internal/models"

	"gorm.io/gorm"
)

type VideoRepository interface {
	Create(video *models.Video) error
	GetByID(id uint, includePending bool) (*models.Video, error) // includePending - показывать ли видео на модерации
	GetFeed(limit, offset int, includePending bool) ([]models.Video, error) // includePending - показывать ли видео на модерации
	GetByUserID(userID uint, limit, offset int) ([]models.Video, error)
	GetByCategory(categoryID uint, limit, offset int, includePending bool) ([]models.Video, error)
	GetByTag(tag string, limit, offset int, includePending bool) ([]models.Video, error)
	GetPendingModeration(limit, offset int) ([]models.Video, error) // Получить видео на модерации
	Update(video *models.Video) error
	Delete(id uint) error
	IncrementViews(id uint) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(video *models.Video) error {
	return r.db.Create(video).Error
}

func (r *videoRepository) GetByID(id uint, includePending bool) (*models.Video, error) {
	var video models.Video
	query := r.db.Preload("User").Preload("Category").Preload("Comments.User").Preload("Likes")
	
	if !includePending {
		query = query.Where("moderation_status = ?", models.ModerationStatusApproved)
	}
	
	err := query.First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) GetFeed(limit, offset int, includePending bool) ([]models.Video, error) {
	var videos []models.Video
	query := r.db.Preload("User").Preload("Category").Order("created_at DESC")
	
	if !includePending {
		query = query.Where("moderation_status = ?", models.ModerationStatusApproved)
	}
	
	err := query.Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) GetPendingModeration(limit, offset int) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Where("moderation_status = ?", models.ModerationStatusPending).
		Preload("User").Preload("Category").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) GetByCategory(categoryID uint, limit, offset int, includePending bool) ([]models.Video, error) {
	var videos []models.Video
	query := r.db.Where("category_id = ?", categoryID).Preload("User").Preload("Category")
	
	if !includePending {
		query = query.Where("moderation_status = ?", models.ModerationStatusApproved)
	}
	
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) GetByTag(tag string, limit, offset int, includePending bool) ([]models.Video, error) {
	var videos []models.Video
	query := r.db.Where("tags LIKE ?", "%"+tag+"%").Preload("User").Preload("Category")
	
	if !includePending {
		query = query.Where("moderation_status = ?", models.ModerationStatusApproved)
	}
	
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) GetByUserID(userID uint, limit, offset int) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Where("user_id = ?", userID).Preload("User").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) Update(video *models.Video) error {
	return r.db.Save(video).Error
}

func (r *videoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Video{}, id).Error
}

func (r *videoRepository) IncrementViews(id uint) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).
		Update("views", gorm.Expr("views + 1")).Error
}
