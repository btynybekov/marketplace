package categories

import (
	"html/template"
	"net/http"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/btynybekov/marketplace/internal/repository"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	Repo repository.Repository
	Tmpl *template.Template
}

func NewCategoriesHandler(repo repository.Repository, tmpl *template.Template) *CategoryHandler {
	return &CategoryHandler{
		Repo: repo,
		Tmpl: tmpl,
	}
}
func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	// Получаем категорию
	category, err := h.Repo.GetCategoryBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// Получаем товары этой категории
	items, err := h.Repo.GetItemsByCategoryID(r.Context(), category.ID)
	if err != nil {
		http.Error(w, "Failed to get items", http.StatusInternalServerError)
		return
	}

	// Формируем структуру для шаблона
	data := struct {
		Category models.Category
		Items    []models.Item
	}{
		Category: category,
		Items:    items,
	}

	// Рендерим шаблон
	if err := h.Tmpl.ExecuteTemplate(w, "category.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
