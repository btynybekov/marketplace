package repository

import (
	"database/sql"

	"github.com/btynybekov/marketplace/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(cat *models.Category) error {
	_, err := r.db.Exec("INSERT INTO categories (name, parent_id) VALUES ($1, $2)", cat.Name, cat.ParentID)
	return err
}

func (r *CategoryRepository) List() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, parent_id FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}
