package homepage

import (
	"html/template"
	"net/http"

	"github.com/btynybekov/marketplace/internal/handlers/shared"
	"github.com/btynybekov/marketplace/internal/repository"
)

type HomePageHandler struct {
	repos repository.RepositorySet
	tmpl  *template.Template
}

func NewHomePageHandler(repos repository.RepositorySet, tmpl *template.Template) *HomePageHandler {
	return &HomePageHandler{repos: repos, tmpl: tmpl}
}

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	roots, err := h.repos.Categories().ListRoots(ctx)
	if err != nil {
		shared.InternalError(w, err)
		return
	}

	if h.tmpl == nil || h.tmpl.Lookup("homepage.html") == nil {
		shared.WriteJSON(w, http.StatusOK, map[string]any{"categories": roots})
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, "homepage.html", map[string]any{
		"Categories": roots,
	}); err != nil {
		http.Error(w, "template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
