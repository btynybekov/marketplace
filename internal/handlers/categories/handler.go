package categories

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/btynybekov/marketplace/internal/handlers/shared"
	"github.com/btynybekov/marketplace/internal/repository"
)

type CategoryHandler struct {
	repos repository.RepositorySet
	tmpl  *template.Template // ожидается layout с именем "categories.html"
}

func NewCategoryHandler(repos repository.RepositorySet, tmpl *template.Template) *CategoryHandler {
	return &CategoryHandler{repos: repos, tmpl: tmpl}
}

// ServeHTTP позволяет использовать как r.Handle("/categories", h)
func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleList(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleList — GET /categories[?parent_slug=transport]
func (h *CategoryHandler) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parent := strings.TrimSpace(r.URL.Query().Get("parent_slug"))

	var (
		data any
		err  error
	)
	if parent == "" {
		data, err = h.repos.Categories().ListRoots(ctx)
	} else {
		data, err = h.repos.Categories().ListChildrenBySlug(ctx, parent)
	}
	if err != nil {
		shared.InternalError(w, err)
		return
	}

	// Если клиент просит JSON — отдадим JSON
	if acceptsJSON(r) || h.tmpl == nil || h.tmpl.Lookup("categories.html") == nil {
		shared.WriteJSON(w, http.StatusOK, data)
		return
	}

	// Иначе отрисуем HTML
	if err := h.tmpl.ExecuteTemplate(w, "categories.html", map[string]any{
		"Categories": data,
		"ParentSlug": parent,
	}); err != nil {
		http.Error(w, "template render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func acceptsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(strings.ToLower(accept), "application/json") ||
		strings.Contains(strings.ToLower(accept), "json")
}

// helper (иногда полезен для query-параметров int)
func atoi(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
