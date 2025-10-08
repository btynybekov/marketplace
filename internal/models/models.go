package models

import (
	"time"

	"github.com/google/uuid"
)

// ===== Каталог =====

type Product struct {
	ID           uuid.UUID         `json:"id"`
	Title        string            `json:"title"`
	PriceAmount  float64           `json:"price_amount"`
	CurrencyCode string            `json:"currency_code"`
	Attrs        map[string]string `json:"attrs,omitempty"`
	FilterURL    string            `json:"filter_url,omitempty"`
}

type ProductMedia struct {
	ProductID uuid.UUID `json:"product_id"`
	URL       string    `json:"url"`
	Sort      int32     `json:"sort"`
	Cover     bool      `json:"cover"`
}

type Category struct {
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	Children []Category `json:"children,omitempty"`
}

// ===== Чат / История =====

type Conversation struct {
	ID        uuid.UUID  `json:"id"`
	SessionID string     `json:"session_id"`        // твой внешний идентификатор сессии (для фронта)
	UserID    *uuid.UUID `json:"user_id,omitempty"` // если есть авторизация — можно NULL
	CreatedAt time.Time  `json:"created_at"`
}

type Message struct {
	ID             uuid.UUID         `json:"id"`
	ConversationID uuid.UUID         `json:"conversation_id"`
	Role           string            `json:"role"` // "user" | "assistant"
	Text           string            `json:"text"`
	Meta           map[string]string `json:"meta,omitempty"` // JSONB
	CreatedAt      time.Time         `json:"created_at"`
}

// Лог поисковых запросов — пригодится для аналитики/персонализации
type SearchRequest struct {
	ID        uuid.UUID         `json:"id"`
	SessionID string            `json:"session_id"`
	UserID    *uuid.UUID        `json:"user_id,omitempty"`
	Query     string            `json:"query"`            // нормализованный текст запроса
	Params    map[string]string `json:"params,omitempty"` // распарсенные фильтры/категории/валюта
	CreatedAt time.Time         `json:"created_at"`
}
