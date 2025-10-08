package items

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/btynybekov/marketplace/internal/handlers/shared"
	"github.com/btynybekov/marketplace/internal/repository"
)

type ItemHandler struct {
	repos repository.RepositorySet
	tmpl  *template.Template // ожидается "items.html"
}

func NewItemHandler(repos repository.RepositorySet, tmpl *template.Template) *ItemHandler {
	return &ItemHandler{repos: repos, tmpl: tmpl}
}

// GET /items?category_slug=cars&limit=20&offset=0
func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	category := strings.TrimSpace(r.URL.Query().Get("category_slug"))
	if category == "" {
		shared.BadRequest(w, "category_slug is required")
		return
	}

	limit := parseInt(r.URL.Query().Get("limit"), 20, 1, 50)
	offset := parseInt(r.URL.Query().Get("offset"), 0, 0, 1000000)

	products, err := h.repos.Products().ListByCategorySlug(ctx, category, limit, offset)
	if err != nil {
		shared.InternalError(w, err)
		return
	}

	if acceptsJSON(r) || h.tmpl == nil || h.tmpl.Lookup("items.html") == nil {
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"items":  products,
			"count":  len(products),
			"limit":  limit,
			"offset": offset,
		})
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, "items.html", map[string]any{
		"Items":        products,
		"CategorySlug": category,
	}); err != nil {
		http.Error(w, "template render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func acceptsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(strings.ToLower(accept), "application/json") ||
		strings.Contains(strings.ToLower(accept), "json")
}

func parseInt(s string, def, min, max int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
