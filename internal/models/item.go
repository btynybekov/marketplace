package models

import "time"

type Item struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	Images      []string  `db:"images" json:"images"`
	CategoryID  int64     `db:"category_id" json:"category_id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
