package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/btynybekov/marketplace/internal/models"
)

// ===== RepositorySet (склейка) =====

type pgRepo struct {
	db *pgxpool.Pool

	productsRepo     ProductsRepository
	productMediaRepo ProductMediaRepository
	categoriesRepo   CategoriesRepository

	conversationsRepo  ConversationsRepository
	messagesRepo       MessagesRepository
	searchRequestsRepo SearchRequestsRepository
}

func New(db *pgxpool.Pool) RepositorySet {
	r := &pgRepo{db: db}
	r.productsRepo = &productsRepo{db: db}
	r.productMediaRepo = &productMediaRepo{db: db}
	r.categoriesRepo = &categoriesRepo{db: db}
	r.conversationsRepo = &conversationsRepo{db: db}
	r.messagesRepo = &messagesRepo{db: db}
	r.searchRequestsRepo = &searchRequestsRepo{db: db}
	return r
}

func (r *pgRepo) Products() ProductsRepository             { return r.productsRepo }
func (r *pgRepo) ProductMedia() ProductMediaRepository     { return r.productMediaRepo }
func (r *pgRepo) Categories() CategoriesRepository         { return r.categoriesRepo }
func (r *pgRepo) Conversations() ConversationsRepository   { return r.conversationsRepo }
func (r *pgRepo) Messages() MessagesRepository             { return r.messagesRepo }
func (r *pgRepo) SearchRequests() SearchRequestsRepository { return r.searchRequestsRepo }

// ===== ProductsRepository impl =====

type productsRepo struct{ db *pgxpool.Pool }

func (r *productsRepo) ListByCategorySlug(ctx context.Context, slug string, limit, offset int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.title, p.price_amount, p.currency_code, p.attrs, p.filter_url
		FROM product p
		JOIN category c ON c.id = p.category_id
		WHERE c.slug = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`, slug, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Product, 0, limit)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Title, &p.PriceAmount, &p.CurrencyCode, &p.Attrs, &p.FilterURL); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ===== ProductMediaRepository impl =====

type productMediaRepo struct{ db *pgxpool.Pool }

func (r *productMediaRepo) ListByProductIDs(ctx context.Context, ids []uuid.UUID) ([]models.ProductMedia, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT product_id, url, sort, cover
		FROM product_media
		WHERE product_id = ANY($1)
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.ProductMedia
	for rows.Next() {
		var m models.ProductMedia
		if err := rows.Scan(&m.ProductID, &m.URL, &m.Sort, &m.Cover); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// ===== CategoriesRepository impl =====

type categoriesRepo struct{ db *pgxpool.Pool }

func (r *categoriesRepo) ListRoots(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT name, slug
		FROM category
		WHERE parent_id IS NULL
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Name, &c.Slug); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *categoriesRepo) Tree(ctx context.Context) ([]models.Category, error) {
	type row struct {
		ID       int64
		Name     string
		Slug     string
		ParentID *int64
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, name, slug, parent_id
		FROM category
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	all := make([]row, 0, 256)
	for rows.Next() {
		var rr row
		if err := rows.Scan(&rr.ID, &rr.Name, &rr.Slug, &rr.ParentID); err != nil {
			return nil, err
		}
		all = append(all, rr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	index := map[int64]*models.Category{}
	for _, r := range all {
		index[r.ID] = &models.Category{Name: r.Name, Slug: r.Slug}
	}
	var roots []models.Category
	for _, r := range all {
		if r.ParentID == nil {
			roots = append(roots, *index[r.ID])
		} else {
			parent := index[*r.ParentID]
			parent.Children = append(parent.Children, *index[r.ID])
		}
	}
	return roots, nil
}

// ===== ConversationsRepository impl =====

type conversationsRepo struct{ db *pgxpool.Pool }

func (r *conversationsRepo) GetOrCreateBySession(ctx context.Context, sessionID string, userID *uuid.UUID) (models.Conversation, error) {
	// пробуем найти
	var c models.Conversation
	err := r.db.QueryRow(ctx, `
		SELECT id, session_id, user_id, created_at
		FROM conversation
		WHERE session_id = $1
		LIMIT 1
	`, sessionID).Scan(&c.ID, &c.SessionID, &c.UserID, &c.CreatedAt)

	if err == nil {
		return c, nil
	}

	// создаём новое
	id := uuid.New()
	now := time.Now().UTC()
	err = r.db.QueryRow(ctx, `
		INSERT INTO conversation (id, session_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, session_id, user_id, created_at
	`, id, sessionID, userID, now).Scan(&c.ID, &c.SessionID, &c.UserID, &c.CreatedAt)
	return c, err
}

func (r *conversationsRepo) GetBySession(ctx context.Context, sessionID string) (models.Conversation, error) {
	var c models.Conversation
	err := r.db.QueryRow(ctx, `
		SELECT id, session_id, user_id, created_at
		FROM conversation
		WHERE session_id = $1
		LIMIT 1
	`, sessionID).Scan(&c.ID, &c.SessionID, &c.UserID, &c.CreatedAt)
	return c, err
}

// ===== MessagesRepository impl =====

type messagesRepo struct{ db *pgxpool.Pool }

func (r *messagesRepo) Append(ctx context.Context, conversationID uuid.UUID, role, text string, meta map[string]string) (uuid.UUID, error) {
	id := uuid.New()
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx, `
		INSERT INTO message (id, conversation_id, role, text, meta, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, conversationID, role, text, meta, now)
	return id, err
}

func (r *messagesRepo) ListLast(ctx context.Context, conversationID uuid.UUID, limit int) ([]models.Message, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, conversation_id, role, text, meta, created_at
		FROM message
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Message, 0, limit)
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Text, &m.Meta, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	// лучше вернуть в порядке по возрастанию времени (UI): разворачиваем
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, rows.Err()
}

// ===== SearchRequestsRepository impl =====

type searchRequestsRepo struct{ db *pgxpool.Pool }

func (r *searchRequestsRepo) Insert(ctx context.Context, sr models.SearchRequest) (uuid.UUID, error) {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	if sr.CreatedAt.IsZero() {
		sr.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO search_request (id, session_id, user_id, query, params, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sr.ID, sr.SessionID, sr.UserID, sr.Query, sr.Params, sr.CreatedAt)
	return sr.ID, err
}
