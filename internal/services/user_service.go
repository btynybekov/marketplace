package services

import (
	"errors"
	"marketplace/internal/models"
	"marketplace/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user *models.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}
