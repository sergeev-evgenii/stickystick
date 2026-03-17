package service

import (
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/repository"
	"time"
)

type AnalyticsService interface {
	Log(ip string, userID *uint, action, resourceType string, resourceID uint, userAgent string) error
	GetStats(since time.Time) (*AnalyticsStats, error)
	GetRecentActivity(limit, offset int) ([]models.ActivityLog, error)
}

type AnalyticsStats struct {
	UniqueVisitors        int64 `json:"unique_visitors"`
	TotalVideoViews       int64 `json:"total_video_views"`
	TotalLogins           int64 `json:"total_logins"`
	TotalRegistrations    int64 `json:"total_registrations"`
	TotalLikes            int64 `json:"total_likes"`
	TotalUploads          int64 `json:"total_uploads"`
	TotalGenerateVideoClicks int64 `json:"total_generate_video_clicks"`
	VideosTotalViews      int64 `json:"videos_total_views"` // сумма views по всем видео (из таблицы videos)
}

type analyticsService struct {
	activityRepo repository.ActivityLogRepository
	videoRepo    repository.VideoRepository
}

func NewAnalyticsService(activityRepo repository.ActivityLogRepository, videoRepo repository.VideoRepository) AnalyticsService {
	return &analyticsService{
		activityRepo: activityRepo,
		videoRepo:    videoRepo,
	}
}

func (s *analyticsService) Log(ip string, userID *uint, action, resourceType string, resourceID uint, userAgent string) error {
	log := &models.ActivityLog{
		IP:           ip,
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserAgent:    userAgent,
	}
	return s.activityRepo.Create(log)
}

func (s *analyticsService) GetStats(since time.Time) (*AnalyticsStats, error) {
	uniqueIPs, _ := s.activityRepo.GetUniqueIPsCount(since)
	videoViews, _ := s.activityRepo.GetTotalViewsCount(since)
	logins, _ := s.activityRepo.GetActionCount(since, models.ActionLogin)
	regs, _ := s.activityRepo.GetActionCount(since, models.ActionRegister)
	likes, _ := s.activityRepo.GetActionCount(since, models.ActionLike)
	uploads, _ := s.activityRepo.GetActionCount(since, models.ActionUpload)
	generateClicks, _ := s.activityRepo.GetActionCount(since, models.ActionGenerateVideoClick)

	// Сумма просмотров по всем видео (из БД) — опционально, можно считать отдельным запросом
	var videosTotalViews int64
	// Упрощённо: не делаем отдельный запрос, можно добавить VideoRepository.SumViews()
	_ = videosTotalViews

	stats := &AnalyticsStats{
		UniqueVisitors:           uniqueIPs,
		TotalVideoViews:          videoViews,
		TotalLogins:              logins,
		TotalRegistrations:       regs,
		TotalLikes:               likes,
		TotalUploads:             uploads,
		TotalGenerateVideoClicks: generateClicks,
	}
	return stats, nil
}

func (s *analyticsService) GetRecentActivity(limit, offset int) ([]models.ActivityLog, error) {
	return s.activityRepo.GetRecent(limit, offset)
}
