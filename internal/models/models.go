package models

import "time"

type Item struct {
	ID            int32
	CategoryId    int32
	SubCategoryId int32
	Stock         *int32
	Price         float64
	Weight        *float64
	Name          string
	Description   string
	ImageURL      *string
	Color         *string
	CreatedAt     time.Time
}

type Category struct {
	ID        int32
	Name      string
	ImageURL  string
	Slug      string
	CreatedAt time.Time
}

type ChatMessage struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type ChatResponse struct {
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
