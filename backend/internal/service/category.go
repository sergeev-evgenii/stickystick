package service

import (
	"sticky-stick/backend/internal/models"
	"sticky-stick/backend/internal/repository"
	"strings"
)

type CategoryService interface {
	GetAll() ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
	GetBySlug(slug string) (*models.Category, error)
	Create(name, slug string) (*models.Category, error)
	Update(id uint, name, slug string) (*models.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetBySlug(slug string) (*models.Category, error) {
	return s.categoryRepo.GetBySlug(slug)
}

func (s *categoryService) Create(name, slug string) (*models.Category, error) {
	// Генерируем slug из name, если не указан
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}

	category := &models.Category{
		Name: name,
		Slug: slug,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Update(id uint, name, slug string) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		category.Name = name
	}
	if slug != "" {
		category.Slug = slug
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id uint) error {
	return s.categoryRepo.Delete(id)
}
