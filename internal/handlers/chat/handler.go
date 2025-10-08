package chat

import (
	"html/template"
	"net/http"

	"github.com/btynybekov/marketplace/internal/handlers/shared"
	"github.com/btynybekov/marketplace/internal/repository"
)

// ChatPageHandler — хендлер страницы чата (GET /chat).
type ChatPageHandler struct {
	repos repository.RepositorySet
	tmpl  *template.Template // ожидается шаблон "chat.html" (если используешь HTML)
}

// NewChatHandler — конструктор страницы чата.
func NewChatHandler(repos repository.RepositorySet) *ChatPageHandler {
	return &ChatPageHandler{repos: repos}
}

// WithTemplate — опционально подвязать html/template.
func (h *ChatPageHandler) WithTemplate(t *template.Template) *ChatPageHandler {
	h.tmpl = t
	return h
}

// ServeHTTP — рендер страницы чата (или JSON-заглушки).
func (h *ChatPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// если нет шаблона — отдадим JSON
	if h.tmpl == nil || h.tmpl.Lookup("chat.html") == nil {
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "chat page",
		})
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, "chat.html", nil); err != nil {
		http.Error(w, "template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
