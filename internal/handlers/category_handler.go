package handlers

import (
	"log"
	"net/http"

	"github.com/btynybekov/marketplace/internal/models"
)

// CategoryPage показывает товары конкретной категории
func (h *Handler) CategoryPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем slug категории из query параметра
	slug := r.URL.Query().Get("category")
	if slug == "" {
		http.Error(w, "категория не указана", http.StatusBadRequest)
		return
	}

	// Получаем категорию по slug
	category, err := h.Repo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		http.Error(w, "категория не найдена", http.StatusNotFound)
		log.Printf("ошибка при получении категории: %v", err)
		return
	}

	// Получаем товары для категории
	items, err := h.Repo.GetItemsByCategoryID(ctx, category.ID)
	if err != nil {
		http.Error(w, "не удалось получить товары", http.StatusInternalServerError)
		log.Printf("ошибка при получении товаров: %v", err)
		return
	}

	data := struct {
		Category models.Category
		Items    []models.Item
	}{
		Category: category,
		Items:    items,
	}

	// Рендерим шаблон category.html внутри layout.html
	if err := h.Tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "ошибка рендеринга шаблона", http.StatusInternalServerError)
		log.Printf("ошибка при выполнении шаблона: %v", err)
	}
}
