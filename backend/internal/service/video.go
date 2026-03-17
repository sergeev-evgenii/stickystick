package service

import (
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/repository"
)

type VideoService interface {
	GetFeed(limit, offset int, isAdmin bool, excludeIDs []uint, orderRandom bool) ([]models.Video, error)
	GetVideo(id uint, isAdmin bool) (*models.Video, error)
	GetByCategory(categoryID uint, limit, offset int, isAdmin bool) ([]models.Video, error)
	GetByTag(tag string, limit, offset int, isAdmin bool) ([]models.Video, error)
	GetPendingModeration(limit, offset int) ([]models.Video, error)
	GetApproved(limit, offset int) ([]models.Video, error)
	GetHidden(limit, offset int) ([]models.Video, error)
	HideVideo(id uint) error
	UnhideVideo(id uint) error
	UpdateVideoFields(id uint, title, description, tags string) error
	UploadVideo(userID uint, title, description, videoURL, thumbnailURL string, duration int) (*models.Video, error)
	UploadMedia(userID uint, title, description, tags string, categoryID *uint, mediaURL string, mediaType models.MediaType, thumbnailURL string, duration int) (*models.Video, error)
	ModerateVideo(videoID uint, status models.ModerationStatus) error
	DeleteVideo(id, userID uint) error
	LikeVideo(videoID, userID uint) error
	UnlikeVideo(videoID, userID uint) error
	AddComment(videoID, userID uint, content string) (*models.Comment, error)
}

type videoService struct {
	videoRepo   repository.VideoRepository
	commentRepo repository.CommentRepository
	likeRepo    repository.LikeRepository
}

func NewVideoService(
	videoRepo repository.VideoRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
) VideoService {
	return &videoService{
		videoRepo:   videoRepo,
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
	}
}

func (s *videoService) GetFeed(limit, offset int, isAdmin bool, excludeIDs []uint, orderRandom bool) ([]models.Video, error) {
	return s.videoRepo.GetFeed(limit, offset, isAdmin, excludeIDs, orderRandom)
}

func (s *videoService) GetByCategory(categoryID uint, limit, offset int, isAdmin bool) ([]models.Video, error) {
	return s.videoRepo.GetByCategory(categoryID, limit, offset, isAdmin)
}

func (s *videoService) GetByTag(tag string, limit, offset int, isAdmin bool) ([]models.Video, error) {
	return s.videoRepo.GetByTag(tag, limit, offset, isAdmin)
}

func (s *videoService) GetPendingModeration(limit, offset int) ([]models.Video, error) {
	return s.videoRepo.GetPendingModeration(limit, offset)
}

func (s *videoService) GetApproved(limit, offset int) ([]models.Video, error) {
	return s.videoRepo.GetApproved(limit, offset)
}

func (s *videoService) GetHidden(limit, offset int) ([]models.Video, error) {
	return s.videoRepo.GetHidden(limit, offset)
}

func (s *videoService) HideVideo(id uint) error {
	return s.videoRepo.SetHidden(id, true)
}

func (s *videoService) UnhideVideo(id uint) error {
	return s.videoRepo.SetHidden(id, false)
}

func (s *videoService) UpdateVideoFields(id uint, title, description, tags string) error {
	return s.videoRepo.UpdateFields(id, title, description, tags)
}

func (s *videoService) GetVideo(id uint, isAdmin bool) (*models.Video, error) {
	video, err := s.videoRepo.GetByID(id, isAdmin)
	if err != nil {
		return nil, err
	}

	// Increment views
	_ = s.videoRepo.IncrementViews(id)

	return video, nil
}

func (s *videoService) ModerateVideo(videoID uint, status models.ModerationStatus) error {
	video, err := s.videoRepo.GetByID(videoID, true) // Админ может видеть все
	if err != nil {
		return err
	}
	
	video.ModerationStatus = status
	return s.videoRepo.Update(video)
}

func (s *videoService) UploadVideo(userID uint, title, description, videoURL, thumbnailURL string, duration int) (*models.Video, error) {
	video := &models.Video{
		UserID:           userID,
		Title:            title,
		Description:      description,
		MediaType:        models.MediaTypeVideo,
		MediaURL:         videoURL,
		ThumbnailURL:     thumbnailURL,
		Duration:         duration,
		ModerationStatus: models.ModerationStatusPending, // Автоматически ставим статус "на модерации"
	}

	if err := s.videoRepo.Create(video); err != nil {
		return nil, err
	}

	return s.videoRepo.GetByID(video.ID, true) // При создании возвращаем с pending
}

func (s *videoService) UploadMedia(userID uint, title, description, tags string, categoryID *uint, mediaURL string, mediaType models.MediaType, thumbnailURL string, duration int) (*models.Video, error) {
	video := &models.Video{
		UserID:           userID,
		CategoryID:       categoryID,
		Title:            title,
		Description:      description,
		Tags:             tags,
		MediaType:        mediaType,
		MediaURL:         mediaURL,
		ThumbnailURL:     thumbnailURL,
		Duration:         duration,
		ModerationStatus: models.ModerationStatusPending, // Автоматически ставим статус "на модерации"
	}

	if err := s.videoRepo.Create(video); err != nil {
		return nil, err
	}

	return s.videoRepo.GetByID(video.ID, true) // При создании возвращаем с pending
}

func (s *videoService) DeleteVideo(id, userID uint) error {
	video, err := s.videoRepo.GetByID(id, true) // Проверяем права, поэтому нужен доступ ко всем
	if err != nil {
		return err
	}

	if video.UserID != userID {
		return repository.ErrUnauthorized
	}

	return s.videoRepo.Delete(id)
}

func (s *videoService) LikeVideo(videoID, userID uint) error {
	exists, err := s.likeRepo.Exists(videoID, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already liked
	}

	like := &models.Like{
		VideoID: videoID,
		UserID:  userID,
	}

	return s.likeRepo.Create(like)
}

func (s *videoService) UnlikeVideo(videoID, userID uint) error {
	return s.likeRepo.Delete(videoID, userID)
}

func (s *videoService) AddComment(videoID, userID uint, content string) (*models.Comment, error) {
	comment := &models.Comment{
		VideoID: videoID,
		UserID:  userID,
		Content: content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// Get comment with user
	comments, err := s.commentRepo.GetByVideoID(videoID, 1, 0)
	if err != nil || len(comments) == 0 {
		return comment, err
	}

	return &comments[0], nil
}
