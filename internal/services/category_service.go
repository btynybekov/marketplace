package services

import (
	"marketplace/internal/models"
	"marketplace/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(cat *models.Category) error {
	return s.repo.Create(cat)
}

func (s *CategoryService) List() ([]models.Category, error) {
	return s.repo.List()
}
