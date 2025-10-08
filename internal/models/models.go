// internal/domain/types.go
package models

import (
	"encoding/json"
	"time"
)

type UUID = string // или pgtype.UUID если используешь pgx-types

type Category struct {
	ID        UUID      `json:"id"`
	ParentID  *UUID     `json:"parent_id,omitempty"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Path      string    `json:"path"`
	IsActive  bool      `json:"is_active"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductMedia struct {
	ID        UUID    `json:"id"`
	ProductID UUID    `json:"product_id"`
	URL       string  `json:"url"`
	Type      string  `json:"type"` // "image" | "video" ...
	IsCover   bool    `json:"is_cover"`
	Alt       *string `json:"alt,omitempty"`
	SortOrder int     `json:"sort_order"`
}

type Product struct {
	ID         UUID            `json:"id"`
	CategoryID UUID            `json:"category_id"`
	BrandID    *UUID           `json:"brand_id,omitempty"`
	Model      string          `json:"model"`
	Title      string          `json:"title"`
	Specs      json.RawMessage `json:"specs"` // JSONB
	IsActive   bool            `json:"is_active"`
	Media      []ProductMedia  `json:"media,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type ListingMedia struct {
	ID        UUID    `json:"id"`
	ListingID UUID    `json:"listing_id"`
	URL       string  `json:"url"`
	Type      string  `json:"type"`
	IsCover   bool    `json:"is_cover"`
	Alt       *string `json:"alt,omitempty"`
	SortOrder int     `json:"sort_order"`
}

type Listing struct {
	ID           UUID            `json:"id"`
	SellerID     UUID            `json:"seller_id"`
	ProductID    UUID            `json:"product_id"`
	CategoryID   UUID            `json:"category_id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	PriceAmount  float64         `json:"price_amount"`
	CurrencyCode string          `json:"currency_code"` // "KGS"/"USD"
	Condition    string          `json:"condition"`     // "new"/"used"
	LocationText string          `json:"location_text"`
	Attrs        json.RawMessage `json:"attrs"`  // JSONB
	Status       string          `json:"status"` // "active"/"archived"
	Media        []ListingMedia  `json:"media,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
