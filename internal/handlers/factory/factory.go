package factory

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/btynybekov/marketplace/config"
	"github.com/btynybekov/marketplace/internal/ai"
	"github.com/btynybekov/marketplace/internal/repository"

	// твои хендлеры
	"github.com/btynybekov/marketplace/internal/handlers/assistant"
	"github.com/btynybekov/marketplace/internal/handlers/categories"
	"github.com/btynybekov/marketplace/internal/handlers/chat"
	"github.com/btynybekov/marketplace/internal/handlers/homepage"
	"github.com/btynybekov/marketplace/internal/handlers/items"
)

type HandlersFactory struct {
	// страницы
	HomepageHandler   *homepage.HomePageHandler
	CategoriesHandler *categories.CategoryHandler
	ItemsHandler      *items.ItemHandler

	// чат
	ChatPageHandler http.Handler      // GET /chat (рендер страницы)
	ChatHandler     *chat.ChatHandler // API: /chat/session, /chat/ajax, /chat/history

	// ассистенты (прямая прокся на n8n)
	BuyerAssistant  *assistant.AssistantHandler
	SellerAssistant *assistant.AssistantHandler
}

// NewHandlersFactory — теперь принимает aiClient.
func NewHandlersFactory(
	repo repository.RepositorySet,
	tmpl *template.Template,
	conf config.EnvConfig,
	aiClient ai.Client,
) *HandlersFactory {

	chatSvc := chat.NewService(repo, aiClient, nil, conf)

	return &HandlersFactory{
		HomepageHandler:   homepage.NewHomePageHandler(repo, tmpl),
		CategoriesHandler: categories.NewCategoryHandler(repo, tmpl),
		ItemsHandler:      items.NewItemHandler(repo, tmpl),

		ChatPageHandler: chat.NewChatHandler(repo).WithTemplate(tmpl), // страница
		ChatHandler:     chat.NewChatHTTP(chatSvc),                    // API-обработчик с методами

		BuyerAssistant:  assistant.NewAssistantHandler(conf.N8NBuyerWebhookURL),
		SellerAssistant: assistant.NewAssistantHandler(conf.N8NSellerWebhookURL),
	}
}

// RegisterRoutes — централизованная регистрация.
func (f *HandlersFactory) RegisterRoutes(r *mux.Router) {
	// ассистенты (n8n)
	r.Handle("/assistant/buyer", f.BuyerAssistant).Methods(http.MethodPost)
	r.Handle("/assistant/seller", f.SellerAssistant).Methods(http.MethodPost)

	// страницы
	r.Handle("/", f.HomepageHandler).Methods(http.MethodGet)
	r.Handle("/chat", f.ChatPageHandler).Methods(http.MethodGet)
	r.Handle("/categories", f.CategoriesHandler).Methods(http.MethodGet)
	r.Handle("/items", f.ItemsHandler).Methods(http.MethodGet)

	// чат API (методами ChatHandler)
	r.Handle("/chat/session", f.ChatHandler.StartSession()).Methods(http.MethodPost)
	r.Handle("/chat/ajax", f.ChatHandler.SendMessage()).Methods(http.MethodPost)
	r.Handle("/chat/history", f.ChatHandler.GetHistory()).Methods(http.MethodGet)
}
