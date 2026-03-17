package repository

import (
	"sticky-stick/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	Create(log *models.ActivityLog) error
	GetRecent(limit, offset int) ([]models.ActivityLog, error)
	GetUniqueIPsCount(since time.Time) (int64, error)
	GetActionCount(since time.Time, action string) (int64, error)
	GetTotalViewsCount(since time.Time) (int64, error)
}

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

func (r *activityLogRepository) GetRecent(limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}

func (r *activityLogRepository) GetUniqueIPsCount(since time.Time) (int64, error) {
	var count int64
	err := r.db.Raw("SELECT COUNT(DISTINCT ip) FROM activity_logs WHERE created_at >= ? AND deleted_at IS NULL", since).Scan(&count).Error
	return count, err
}

func (r *activityLogRepository) GetActionCount(since time.Time, action string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActivityLog{}).Where("created_at >= ? AND action = ?", since, action).Count(&count).Error
	return count, err
}

func (r *activityLogRepository) GetTotalViewsCount(since time.Time) (int64, error) {
	return r.GetActionCount(since, models.ActionVideoView)
}
