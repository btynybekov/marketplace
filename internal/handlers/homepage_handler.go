package handlers

import (
	"log"
	"net/http"

	"github.com/btynybekov/marketplace/internal/models"
)

type HomePageData struct {
	Categories  []models.Category
	RecentItems []models.Item
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем категории
	categories, err := h.Repo.GetCategories(ctx)
	if err != nil {
		http.Error(w, "не удалось получить категории", http.StatusInternalServerError)
		log.Printf("ошибка при получении категорий: %v", err)
		return
	}

	// Получаем последние товары
	items, err := h.Repo.GetRecentItems(ctx, 10)
	if err != nil {
		http.Error(w, "не удалось получить товары", http.StatusInternalServerError)
		log.Printf("ошибка при получении товаров: %v", err)
		return
	}

	data := HomePageData{
		Categories:  categories,
		RecentItems: items,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "ошибка рендеринга шаблона", http.StatusInternalServerError)
		log.Printf("ошибка при выполнении шаблона: %v", err)
	}
}
