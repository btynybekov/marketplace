package repository

import (
	"context"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	DB *pgxpool.Pool
}

var _ Repository = (*PostgresRepo)(nil)

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

func (r *PostgresRepo) GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, name, created_at, image_url, slug FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.CreatedAt,
			&c.ImageURL,
			&c.Slug,
		); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *PostgresRepo) GetRecentItems(ctx context.Context, limit int) ([]models.Item, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, name, description, price, category_id, created_at, sub_category_id, stock, weight, color, image_url
        FROM products
        ORDER BY created_at DESC
        LIMIT $1
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.CategoryId,
			&item.CreatedAt,
			&item.SubCategoryId,
			&item.Stock,
			&item.Weight,
			&item.Color,
			&item.ImageURL,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *PostgresRepo) GetCategoryBySlug(ctx context.Context, slug string) (models.Category, error) {
	var c models.Category
	query := `SELECT id, name, slug, image_url FROM categories WHERE slug=$1`
	err := r.DB.QueryRow(ctx, query, slug).Scan(
		&c.ID,
		&c.Name,
		&c.Slug,
		&c.ImageURL,
	)
	return c, err
}

func (r *PostgresRepo) GetItemsByCategoryID(ctx context.Context, categoryID int32) ([]models.Item, error) {
	var items []models.Item
	query := `SELECT id, category_id, name, price, description, image_url FROM products WHERE category_id=$1`
	rows, err := r.DB.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.CategoryId,
			&item.Name,
			&item.Price,
			&item.Description,
			&item.ImageURL,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *PostgresRepo) SaveMessage(ctx context.Context, userID, role, content string) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO chat_messages (user_id, role, content, created_at)
		VALUES ($1, $2, $3, NOW())
	`, userID, role, content)
	return err
}

func (r *PostgresRepo) GetHistory(ctx context.Context, userID string, limit int) ([]models.ChatMessage, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, user_id, role, content, created_at
		FROM chat_messages
		WHERE user_id = $1
		ORDER BY created_at ASC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		if err := rows.Scan(&m.ID, &m.UserID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *PostgresRepo) GetItemByID(ctx context.Context, itemID int32) (models.Item, error) {
	var item models.Item
	query := `
		SELECT id, name, description, price, category_id, created_at, sub_category_id, stock, weight, color, image_url
		FROM products
		WHERE id = $1
	`
	err := r.DB.QueryRow(ctx, query, itemID).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.CategoryId,
		&item.CreatedAt,
		&item.SubCategoryId,
		&item.Stock,
		&item.Weight,
		&item.Color,
		&item.ImageURL,
	)
	return item, err
}
