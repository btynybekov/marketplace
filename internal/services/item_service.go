package services

import (
	"errors"
	"marketplace/internal/models"
	"marketplace/internal/repository"
)

var ErrForbidden = errors.New("forbidden")

type ItemService struct {
	repo         *repository.ItemRepository
	categoryRepo *repository.CategoryRepository
}

func NewItemService(repo *repository.ItemRepository, categoryRepo *repository.CategoryRepository) *ItemService {
	return &ItemService{repo: repo, categoryRepo: categoryRepo}
}

func (s *ItemService) Create(item *models.Item) error {
	return s.repo.Create(item)
}

func (s *ItemService) GetByID(id int64) (*models.Item, error) {
	return s.repo.GetByID(id)
}

func (s *ItemService) Update(item *models.Item, userID int64) error {
	existing, err := s.repo.GetByID(item.ID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Update(item)
}

func (s *ItemService) Delete(id int64, userID int64) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

func (s *ItemService) List() ([]models.Item, error) {
	return s.repo.List()
}
