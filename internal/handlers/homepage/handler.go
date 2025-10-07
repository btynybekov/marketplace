// internal/handlers/homepage/homepage.go
package homepage

import (
	"html/template"
	"net/http"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/btynybekov/marketplace/internal/repository"
)

type HomePageHandler struct {
	Repo repository.Repository
	Tmpl *template.Template
}

func NewHomePageHandler(repo repository.Repository, tmpl *template.Template) *HomePageHandler {
	return &HomePageHandler{
		Repo: repo,
		Tmpl: tmpl,
	}
}

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Categories  []models.Category
		RecentItems []models.Item
	}{}

	data.Categories, _ = h.Repo.GetCategories(r.Context())
	data.RecentItems, _ = h.Repo.GetRecentItems(r.Context(), 10)

	if err := h.Tmpl.ExecuteTemplate(w, "homepage.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
