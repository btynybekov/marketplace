package items

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/btynybekov/marketplace/internal/repository"
)

type ItemHandler struct {
	Repo repository.Repository
	Tmpl *template.Template
}

func NewItemsHandler(repo repository.Repository, tmpl *template.Template) *ItemHandler {
	return &ItemHandler{
		Repo: repo,
		Tmpl: tmpl,
	}
}

func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := h.Repo.GetItemByID(r.Context(), int32(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.Tmpl.ExecuteTemplate(w, "items.html", item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
