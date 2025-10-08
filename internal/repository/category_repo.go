package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/btynybekov/marketplace/storage"
)

//
// ──────────────────────────────────────────────
//   ИНТЕРФЕЙСЫ
// ──────────────────────────────────────────────
//

// CategoryRepository управляет категориями.
type CategoryRepository interface {
	ListRoots(ctx context.Context) ([]models.Category, error)
	ListChildrenBySlug(ctx context.Context, parentSlug string) ([]models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
}

// ProductRepository управляет товарами.
type ProductRepository interface {
	Get(ctx context.Context, id string) (*models.Product, error)
	ListByCategorySlug(ctx context.Context, slug string, limit, offset int) ([]models.Product, error)
}

// ListingRepository управляет объявлениями.
type ListingRepository interface {
	Get(ctx context.Context, id string) (*models.Listing, error)
	Search(ctx context.Context, p ListingSearchParams) ([]ListingWithProduct, error)
}

// RepositorySet агрегирует все репозитории.
type RepositorySet interface {
	Categories() CategoryRepository
	Products() ProductRepository
	Listings() ListingRepository
}

// ListingSearchParams — фильтры для поиска объявлений.
type ListingSearchParams struct {
	CategorySlug string
	PriceMax     *float64
	Currency     string
	Attrs        map[string]string
	SortBy       string
	SortOrder    string
	Limit        int
	Offset       int
}

// ListingWithProduct — результат поиска (объявление + продукт + медиа).
type ListingWithProduct struct {
	Listing models.Listing        `json:"listing"`
	Product models.Product        `json:"product"`
	Media   []models.ListingMedia `json:"media,omitempty"`
}

//
// ──────────────────────────────────────────────
//   РЕАЛИЗАЦИЯ НА POSTGRES
// ──────────────────────────────────────────────
//

// postgresSet агрегирует конкретные реализации.
type postgresSet struct {
	cats  *categoryRepo
	prods *productRepo
	lists *listingRepo
}

// NewSet возвращает готовый набор репозиториев.
func NewSet(db *storage.DB) RepositorySet {
	return &postgresSet{
		cats:  &categoryRepo{db: db},
		prods: &productRepo{db: db},
		lists: &listingRepo{db: db},
	}
}

func (s *postgresSet) Categories() CategoryRepository { return s.cats }
func (s *postgresSet) Products() ProductRepository    { return s.prods }
func (s *postgresSet) Listings() ListingRepository    { return s.lists }

//
// ──────────────────────────────────────────────
//   CATEGORY REPO
// ──────────────────────────────────────────────
//

type categoryRepo struct{ db *storage.DB }

func (r *categoryRepo) ListRoots(ctx context.Context) ([]models.Category, error) {
	const q = `
SELECT id, parent_id, name, slug, path, is_active, sort_order, created_at, updated_at
FROM category
WHERE parent_id IS NULL AND is_active = TRUE
ORDER BY sort_order, name;`
	rows, err := r.db.Pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Name, &c.Slug, &c.Path, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *categoryRepo) ListChildrenBySlug(ctx context.Context, slug string) ([]models.Category, error) {
	const q = `
SELECT c2.id, c2.parent_id, c2.name, c2.slug, c2.path, c2.is_active, c2.sort_order, c2.created_at, c2.updated_at
FROM category c1
JOIN category c2 ON c2.parent_id = c1.id
WHERE c1.slug = $1 AND c2.is_active = TRUE
ORDER BY c2.sort_order, c2.name;`
	rows, err := r.db.Pool.Query(ctx, q, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Name, &c.Slug, &c.Path, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *categoryRepo) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	const q = `
SELECT id, parent_id, name, slug, path, is_active, sort_order, created_at, updated_at
FROM category WHERE slug=$1;`
	var c models.Category
	if err := r.db.Pool.QueryRow(ctx, q, slug).
		Scan(&c.ID, &c.ParentID, &c.Name, &c.Slug, &c.Path, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

//
// ──────────────────────────────────────────────
//   PRODUCT REPO
// ──────────────────────────────────────────────
//

type productRepo struct{ db *storage.DB }

func (r *productRepo) Get(ctx context.Context, id string) (*models.Product, error) {
	const q = `
SELECT id, category_id, brand_id, model, title, specs, is_active, created_at, updated_at
FROM product WHERE id=$1;`
	var p models.Product
	if err := r.db.Pool.QueryRow(ctx, q, id).
		Scan(&p.ID, &p.CategoryID, &p.BrandID, &p.Model, &p.Title, &p.Specs, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}

	const qm = `
SELECT id, product_id, url, type, is_cover, alt, sort_order
FROM product_media
WHERE product_id=$1
ORDER BY is_cover DESC, sort_order ASC;`
	rows, err := r.db.Pool.Query(ctx, qm, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m models.ProductMedia
		if err := rows.Scan(&m.ID, &m.ProductID, &m.URL, &m.Type, &m.IsCover, &m.Alt, &m.SortOrder); err != nil {
			return nil, err
		}
		p.Media = append(p.Media, m)
	}
	return &p, rows.Err()
}

func (r *productRepo) ListByCategorySlug(ctx context.Context, slug string, limit, offset int) ([]models.Product, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	const q = `
SELECT p.id, p.category_id, p.brand_id, p.model, p.title, p.specs, p.is_active, p.created_at, p.updated_at
FROM product p
JOIN category c ON c.id=p.category_id
WHERE c.slug=$1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;`
	rows, err := r.db.Pool.Query(ctx, q, slug, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Product, 0, limit)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.BrandID, &p.Model, &p.Title, &p.Specs, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

//
// ──────────────────────────────────────────────
//   LISTING REPO
// ──────────────────────────────────────────────
//

type listingRepo struct{ db *storage.DB }

func (r *listingRepo) Get(ctx context.Context, id string) (*models.Listing, error) {
	const q = `
SELECT id, seller_id, product_id, category_id, title, description,
       price_amount, currency_code, condition, location_text,
       attrs, status, created_at, updated_at
FROM listing WHERE id=$1;`
	var l models.Listing
	if err := r.db.Pool.QueryRow(ctx, q, id).
		Scan(&l.ID, &l.SellerID, &l.ProductID, &l.CategoryID, &l.Title, &l.Description,
			&l.PriceAmount, &l.CurrencyCode, &l.Condition, &l.LocationText,
			&l.Attrs, &l.Status, &l.CreatedAt, &l.UpdatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *listingRepo) Search(ctx context.Context, p ListingSearchParams) ([]ListingWithProduct, error) {
	limit := p.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := 0
	if p.Offset > 0 {
		offset = p.Offset
	}

	sortBy := "l.price_amount"
	if p.SortBy == "created_at" {
		sortBy = "l.created_at"
	}
	sortOrder := "ASC"
	if strings.EqualFold(p.SortOrder, "desc") {
		sortOrder = "DESC"
	}

	where := []string{"c.slug = $1", "l.status = 'active'", "l.currency_code = $2"}
	args := []any{p.CategorySlug, p.Currency}
	argPos := 3

	if p.PriceMax != nil {
		where = append(where, fmt.Sprintf("l.price_amount <= $%d", argPos))
		args = append(args, *p.PriceMax)
		argPos++
	}

	keys := make([]string, 0, len(p.Attrs))
	for k := range p.Attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := p.Attrs[k]
		where = append(where, fmt.Sprintf("l.attrs->>$%d = $%d", argPos, argPos+1))
		args = append(args, k, v)
		argPos += 2
	}

	q := fmt.Sprintf(`
SELECT
  l.id, l.seller_id, l.product_id, l.category_id, l.title, l.description,
  l.price_amount, l.currency_code, l.condition, l.location_text, l.attrs, l.status, l.created_at, l.updated_at,
  p.id, p.category_id, p.brand_id, p.model, p.title, p.specs, p.is_active, p.created_at, p.updated_at
FROM listing l
JOIN category c ON c.id=l.category_id
JOIN product p  ON p.id=l.product_id
WHERE %s
ORDER BY %s %s, l.created_at DESC
LIMIT %d OFFSET %d;`,
		strings.Join(where, " AND "), sortBy, sortOrder, limit, offset,
	)

	rows, err := r.db.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ListingWithProduct
	var ids []string

	for rows.Next() {
		var li models.Listing
		var pr models.Product
		if err := rows.Scan(
			&li.ID, &li.SellerID, &li.ProductID, &li.CategoryID, &li.Title, &li.Description,
			&li.PriceAmount, &li.CurrencyCode, &li.Condition, &li.LocationText, &li.Attrs, &li.Status, &li.CreatedAt, &li.UpdatedAt,
			&pr.ID, &pr.CategoryID, &pr.BrandID, &pr.Model, &pr.Title, &pr.Specs, &pr.IsActive, &pr.CreatedAt, &pr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, ListingWithProduct{Listing: li, Product: pr})
		ids = append(ids, li.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return results, nil
	}

	qm := `
SELECT listing_id, url, is_cover, sort_order
FROM listing_media
WHERE listing_id = ANY($1)
ORDER BY listing_id, sort_order;`
	mr, err := r.db.Pool.Query(ctx, qm, ids)
	if err != nil {
		return nil, err
	}
	defer mr.Close()

	type mediaRow struct {
		ListingID string
		URL       string
		IsCover   bool
		Sort      int
	}
	mediaMap := make(map[string][]models.ListingMedia)
	for mr.Next() {
		var id string
		var m models.ListingMedia
		if err := mr.Scan(&id, &m.URL, &m.IsCover, &m.SortOrder); err != nil {
			return nil, err
		}
		mediaMap[id] = append(mediaMap[id], m)
	}
	if err := mr.Err(); err != nil {
		return nil, err
	}

	for i := range results {
		if m, ok := mediaMap[results[i].Listing.ID]; ok {
			results[i].Media = m
		}
	}
	return results, nil
}
