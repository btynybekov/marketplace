package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/btynybekov/marketplace/internal/models"
)

// ===== Каталог =====

type ProductsRepository interface {
	ListByCategorySlug(ctx context.Context, slug string, limit, offset int) ([]models.Product, error)
}

type ProductMediaRepository interface {
	ListByProductIDs(ctx context.Context, ids []uuid.UUID) ([]models.ProductMedia, error)
}

type CategoriesRepository interface {
	ListRoots(ctx context.Context) ([]models.Category, error)
	Tree(ctx context.Context) ([]models.Category, error)
}

// ===== Чат / История =====

// Храним/находим разговор по session_id (который фронт держит у себя).
type ConversationsRepository interface {
	GetOrCreateBySession(ctx context.Context, sessionID string, userID *uuid.UUID) (models.Conversation, error)
	GetBySession(ctx context.Context, sessionID string) (models.Conversation, error)
}

// Сообщения внутри разговора.
type MessagesRepository interface {
	Append(ctx context.Context, conversationID uuid.UUID, role, text string, meta map[string]string) (uuid.UUID, error)
	ListLast(ctx context.Context, conversationID uuid.UUID, limit int) ([]models.Message, error)
}

// Лог поисковых запросов (для аналитики/персонализации).
type SearchRequestsRepository interface {
	Insert(ctx context.Context, sr models.SearchRequest) (uuid.UUID, error)
}

// ===== Набор всех репозиториев =====

type RepositorySet interface {
	// каталог
	Products() ProductsRepository
	ProductMedia() ProductMediaRepository
	Categories() CategoriesRepository

	// чат
	Conversations() ConversationsRepository
	Messages() MessagesRepository
	SearchRequests() SearchRequestsRepository
}
