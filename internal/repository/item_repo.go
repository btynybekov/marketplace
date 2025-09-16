package repository

import (
	"database/sql"
	"marketplace/internal/models"
	"time"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(item *models.Item) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		"INSERT INTO items (title, description, price, images, category_id, user_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		item.Title, item.Description, item.Price, item.Images, item.CategoryID, item.UserID, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *ItemRepository) GetByID(id int64) (*models.Item, error) {
	item := &models.Item{}
	err := r.db.QueryRow(
		"SELECT id, title, description, price, images, category_id, user_id, created_at, updated_at FROM items WHERE id=$1", id,
	).Scan(&item.ID, &item.Title, &item.Description, &item.Price, &item.Images, &item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) Update(item *models.Item) error {
	item.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		"UPDATE items SET title=$1, description=$2, price=$3, images=$4, category_id=$5, updated_at=$6 WHERE id=$7",
		item.Title, item.Description, item.Price, item.Images, item.CategoryID, item.UpdatedAt, item.ID,
	)
	return err
}

func (r *ItemRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM items WHERE id=$1", id)
	return err
}

func (r *ItemRepository) List() ([]models.Item, error) {
	rows, err := r.db.Query("SELECT id, title, description, price, images, category_id, user_id, created_at, updated_at FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Price, &item.Images, &item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
