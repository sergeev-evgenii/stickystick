package service

import (
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/repository"
)

type UserService interface {
	GetProfile(id uint) (*models.User, error)
	UpdateProfile(id uint, username, bio, avatar string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateProfile(id uint, username, bio, avatar string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if username != "" {
		user.Username = username
	}
	if bio != "" {
		user.Bio = bio
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
