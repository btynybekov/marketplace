package models

import "time"

type Item struct {
	ID            int32     `db:"id"`
	CategoryId    int32     `db:"category_id"`
	SubCategoryId int32     `db:"sub_category_id"`
	Stock         *int32    `db:"stock"`
	Price         float64   `db:"price"`
	Weight        *float64  `db:"weight"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	ImageURL      *string   `db:"image_url"`
	Color         *string   `db:"color"`
	CreatedAt     time.Time `db:"created_at"`
}

type Category struct {
	ID        int32     `db:"id"`
	Name      string    `db:"name"`
	ImageURL  string    `db:"image_url"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
}
