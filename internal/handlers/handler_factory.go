package handlers

import (
	"log"
	"text/template"

	"github.com/btynybekov/marketplace/internal/repository"
)

type Handler struct {
	Repo repository.Repository
	Tmpl *template.Template
}

// NewHandler создает новый Handler с подключенными шаблонами
func NewHandler(repo repository.Repository) *Handler {
	// Парсим все шаблоны из папки templates
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("ошибка при парсинге шаблонов: %v", err)
	}

	return &Handler{
		Repo: repo,
		Tmpl: tmpl,
	}
}
