package repository

import (
	"database/sql"

	"github.com/btynybekov/marketplace/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash, phone, avatar_url, created_at) VALUES ($1,$2,$3,$4,$5,NOW())",
		user.Name, user.Email, user.PasswordHash, user.Phone, user.AvatarURL,
	)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, phone, avatar_url, created_at FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Phone, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, phone, avatar_url, created_at FROM users WHERE id=$1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Phone, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
