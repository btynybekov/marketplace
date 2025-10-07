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
	GetItemByID(ctx context.Context, itemID int32) (models.Item, error)
	SaveMessage(ctx context.Context, userID, role, message string) error
	GetHistory(ctx context.Context, userID string, limit int) ([]models.ChatMessage, error)
}
