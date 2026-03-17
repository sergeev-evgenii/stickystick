package repository

import (
	"sticky-stick/backend/internal/models"

	"gorm.io/gorm"
)

type VideoReportRepository interface {
	Create(report *models.VideoReport) error
	Exists(videoID uint, userID *uint) (bool, error)
	GetReportedVideoIDs(limit, offset int) ([]uint, error)
	GetReportCount(videoID uint) (int64, error)
}

type videoReportRepository struct {
	db *gorm.DB
}

func NewVideoReportRepository(db *gorm.DB) VideoReportRepository {
	return &videoReportRepository{db: db}
}

func (r *videoReportRepository) Create(report *models.VideoReport) error {
	return r.db.Create(report).Error
}

func (r *videoReportRepository) Exists(videoID uint, userID *uint) (bool, error) {
	if userID == nil {
		return false, nil // анонимные жалобы не дедуплицируем
	}
	var count int64
	err := r.db.Model(&models.VideoReport{}).Where("video_id = ? AND user_id = ?", videoID, *userID).Count(&count).Error
	return count > 0, err
}

// GetReportedVideoIDs возвращает ID видео, на которые есть хотя бы одна жалоба (по убыванию количества жалоб)
func (r *videoReportRepository) GetReportedVideoIDs(limit, offset int) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.VideoReport{}).
		Select("video_id").
		Group("video_id").
		Order("COUNT(*) DESC").
		Limit(limit).
		Offset(offset).
		Pluck("video_id", &ids).Error
	return ids, err
}

func (r *videoReportRepository) GetReportCount(videoID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.VideoReport{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}
