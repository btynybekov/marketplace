package chat

import (
	"html/template"
	"log"
	"net/http"

	"github.com/btynybekov/marketplace/internal/repository"
)

type ChatHandler struct {
	Repo repository.Repository
	Tmpl *template.Template
}

func NewChatHandler(repo repository.Repository) *ChatHandler {
	tmpl, err := template.ParseFiles("templates/chat.html") // только этот шаблон
	if err != nil {
		log.Fatalf("ошибка при парсинге chat.html: %v", err)
	}

	return &ChatHandler{
		Repo: repo,
		Tmpl: tmpl,
	}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "Чат с ассистентом",
	}
	err := h.Tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}
