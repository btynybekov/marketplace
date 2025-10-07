package factory

import (
	"html/template"
	"log"
	"net/http"

	"github.com/btynybekov/marketplace/config"
	"github.com/btynybekov/marketplace/internal/handlers/assistant"
	"github.com/btynybekov/marketplace/internal/handlers/categories"
	"github.com/btynybekov/marketplace/internal/handlers/chat"
	"github.com/btynybekov/marketplace/internal/handlers/homepage"
	"github.com/btynybekov/marketplace/internal/handlers/items"
	"github.com/btynybekov/marketplace/internal/middleware"
	"github.com/btynybekov/marketplace/internal/repository"
	"github.com/gorilla/mux"
)

// HandlersFactory хранит все хендлеры
type HandlersFactory struct {
	ChatHandler       *chat.ChatHandler
	ChatAjaxHandler   *chat.ChatAjaxHandler
	CategoriesHandler *categories.CategoryHandler
	ItemsHandler      *items.ItemHandler
	HomepageHandler   *homepage.HomePageHandler
	BuyerAssistant    *assistant.AssistantHandler
	SellerAssistant   *assistant.AssistantHandler
}

// NewHandlersFactory создаёт все хендлеры и возвращает фабрику
func NewHandlersFactory(repo repository.Repository, tmpl *template.Template, conf config.EnvConfig) *HandlersFactory {
	log.Println("Buyer Webhook URL:", conf.N8NBuyerWebhookURL)
	log.Println("Seller Webhook URL:", conf.N8NSellerWebhookURL)
	return &HandlersFactory{
		ChatHandler:       chat.NewChatHandler(repo),
		ChatAjaxHandler:   chat.NewChatAjaxHandler(repo),
		CategoriesHandler: categories.NewCategoriesHandler(repo, tmpl),
		ItemsHandler:      items.NewItemsHandler(repo, tmpl),
		HomepageHandler:   homepage.NewHomePageHandler(repo, tmpl),
		BuyerAssistant:    assistant.NewAssistantHandler(conf.N8NBuyerWebhookURL),
		SellerAssistant:   assistant.NewAssistantHandler(conf.N8NSellerWebhookURL),
	}
}

// RegisterRoutes регистрирует все маршруты через mux.Router
func (f *HandlersFactory) RegisterRoutes(r *mux.Router) {
	r.Handle("/assistant/buyer", f.BuyerAssistant).Methods(http.MethodPost)
	r.Handle("/assistant/seller", middleware.AuthMiddleware(f.SellerAssistant)).Methods(http.MethodPost)
	r.Handle("/", f.HomepageHandler).Methods(http.MethodGet)
	r.Handle("/chat", f.ChatHandler).Methods(http.MethodGet)
	r.Handle("/chat/ajax", f.ChatAjaxHandler).Methods(http.MethodPost)
	r.Handle("/categories", f.CategoriesHandler).Methods(http.MethodGet)
	r.Handle("/items", f.ItemsHandler).Methods(http.MethodGet)
}

// MustParseTemplates парсит все html шаблоны и логирует fatal при ошибке
func MustParseTemplates(pattern string) *template.Template {
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	return tmpl
}
