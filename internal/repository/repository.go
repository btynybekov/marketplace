package repository

import (
	"context"

	"github.com/btynybekov/marketplace/internal/models"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetRecentItems(ctx context.Context, limit int) ([]models.Item, error)
	GetCategoryBySlug(ctx context.Context, slug string) (models.Category, error)
	GetItemsByCategoryID(ctx context.Context, categoryID int32) ([]models.Item, error)
}
