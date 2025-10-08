package factory

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/btynybekov/marketplace/config"
	"github.com/btynybekov/marketplace/internal/ai"
	"github.com/btynybekov/marketplace/internal/repository"

	"github.com/btynybekov/marketplace/internal/handlers/assistant"
	"github.com/btynybekov/marketplace/internal/handlers/categories"
	"github.com/btynybekov/marketplace/internal/handlers/chat"
	"github.com/btynybekov/marketplace/internal/handlers/homepage"
	"github.com/btynybekov/marketplace/internal/handlers/items"
)

type HandlersFactory struct {
	HomepageHandler   *homepage.HomePageHandler
	CategoriesHandler *categories.CategoryHandler
	ItemsHandler      *items.ItemHandler
	ChatPageHandler   http.Handler
	ChatHandler       *chat.ChatHandler // методы: StartSession, SendMessage, GetHistory

	// ассистенты (проксирование в n8n)
	BuyerAssistant  *assistant.AssistantHandler
	SellerAssistant *assistant.AssistantHandler
}

// NewHandlersFactory — собирает все зависимости и создаёт хендлеры.
func NewHandlersFactory(
	repo repository.RepositorySet,
	tmpl *template.Template,
	conf config.EnvConfig,
	aiClient ai.Client,
) *HandlersFactory {
	// Сервис чата: LLM + авто выбор buyer/seller по контексту
	chatSvc := chat.NewService(repo, aiClient, nil, conf)

	return &HandlersFactory{
		HomepageHandler:   homepage.NewHomePageHandler(repo, tmpl),
		CategoriesHandler: categories.NewCategoryHandler(repo, tmpl),
		ItemsHandler:      items.NewItemHandler(repo, tmpl),
		ChatPageHandler:   chat.NewChatHandler(repo).WithTemplate(tmpl),
		ChatHandler:       chat.NewChatHTTP(chatSvc),
		// Ассистенты (прямые вебхуки n8n)
		BuyerAssistant:  assistant.NewAssistantHandler(conf.N8NBuyerWebhookURL),
		SellerAssistant: assistant.NewAssistantHandler(conf.N8NSellerWebhookURL),
	}
}

// RegisterRoutes — регистрирует все маршруты на переданный *mux.Router.
func (f *HandlersFactory) RegisterRoutes(r *mux.Router) {
	// Страницы
	r.Handle("/", f.HomepageHandler).Methods(http.MethodGet)
	r.Handle("/categories", f.CategoriesHandler).Methods(http.MethodGet)
	r.Handle("/items", f.ItemsHandler).Methods(http.MethodGet)
	r.Handle("/chat", f.ChatPageHandler).Methods(http.MethodGet)
	// API чата
	r.Handle("/chat/session", f.ChatHandler.StartSession()).Methods(http.MethodPost)
	r.Handle("/chat/ajax", f.ChatHandler.SendMessage()).Methods(http.MethodPost)
	r.Handle("/chat/history", f.ChatHandler.GetHistory()).Methods(http.MethodGet)
	// Ассистенты из n8n webhook
	r.Handle("/assistant/buyer", f.BuyerAssistant).Methods(http.MethodPost)
	r.Handle("/assistant/seller", f.SellerAssistant).Methods(http.MethodPost)
}
