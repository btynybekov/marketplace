package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/btynybekov/marketplace/config"
	"github.com/btynybekov/marketplace/internal/ai"
	"github.com/btynybekov/marketplace/internal/repository"
)

// Service — интерфейс, который использует твой ChatHandler (http.go).
type Service interface {
	StartSession(r *http.Request, userID, sessionID string) (string, error)
	AppendUserMessage(r *http.Request, sessionID, text string, meta map[string]string) (string, error)
	GenerateAssistantReply(r *http.Request, sessionID string) (reply string, extra map[string]any, err error)
	GetHistory(r *http.Request, sessionID string, limit int) ([]messageDTO, error)
}

// service — конкретная реализация Service.
type service struct {
	repos      repository.RepositorySet
	ai         ai.Client
	httpClient *http.Client
	cfg        config.EnvConfig
}

// NewService — создаёт новый сервис чата.
func NewService(repos repository.RepositorySet, aiClient ai.Client, httpClient *http.Client, cfg config.EnvConfig) Service {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &service{
		repos:      repos,
		ai:         aiClient,
		httpClient: httpClient,
		cfg:        cfg,
	}
}

// StartSession — создаёт (или возвращает) session_id.
func (s *service) StartSession(_ *http.Request, userID, sessionID string) (string, error) {
	if sessionID != "" {
		return sessionID, nil
	}
	return "sess-" + time.Now().Format("20060102150405"), nil
}

// AppendUserMessage — сохраняет сообщение пользователя (если нужно).
func (s *service) AppendUserMessage(_ *http.Request, sessionID, text string, meta map[string]string) (string, error) {
	if sessionID == "" || text == "" {
		return "", errors.New("session_id and text required")
	}
	// TODO: сохранить в таблицу message (если хочешь вести историю)
	return "msg-" + time.Now().Format("150405"), nil
}

// GetHistory — возвращает историю сообщений.
func (s *service) GetHistory(_ *http.Request, sessionID string, limit int) ([]messageDTO, error) {
	// TODO: извлечь последние N сообщений из БД
	return []messageDTO{}, nil
}

// GenerateAssistantReply — решает, что делать: болталка или buyer/seller.
func (s *service) GenerateAssistantReply(r *http.Request, sessionID string) (string, map[string]any, error) {
	ctx := r.Context()
	userText := extractLastUserText(r)

	// определить намерение (buy/sell/chitchat)
	intent, err := s.classifyIntent(ctx, userText)
	if err != nil {
		intent = "chitchat"
	}

	switch intent {
	case "buy":
		return s.handleBuyer(ctx, sessionID, userText)
	case "sell":
		return s.handleSeller(ctx, sessionID, userText)
	default:
		// просто болталка
		reply, err := s.ai.Chat(ctx, s.cfg.AIModel, s.cfg.AITemperature, []ai.Message{
			{Role: "system", Content: "Ты дружелюбный помощник маркетплейса."},
			{Role: "user", Content: userText},
		})
		if err != nil {
			return "", nil, err
		}
		return reply, nil, nil
	}
}

//
// ─── PRIVATE HELPERS ───────────────────────────────────────────────────────────
//

func (s *service) classifyIntent(ctx context.Context, text string) (string, error) {
	prompt := `Определи намерение пользователя как одно слово из списка [buy, sell, chitchat].
Текст: ` + text

	resp, err := s.ai.Chat(ctx, s.cfg.AIModel, 0.0, []ai.Message{
		{Role: "system", Content: "Ты классификатор намерений."},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return "chitchat", err
	}

	resp = strings.ToLower(resp)
	switch {
	case strings.Contains(resp, "buy"):
		return "buy", nil
	case strings.Contains(resp, "sell"):
		return "sell", nil
	default:
		return "chitchat", nil
	}
}

// handleBuyer — вызывает buyer webhook (n8n или твой /api/search).
func (s *service) handleBuyer(ctx context.Context, sessionID, text string) (string, map[string]any, error) {
	url := s.cfg.N8NBuyerWebhookURL
	if url == "" {
		return "Buyer webhook не настроен", nil, nil
	}

	payload := map[string]any{
		"intent": "search_listings",
		"requirements": map[string]any{
			"category_slug": "phones",
			"price_max":     15000,
			"currency":      "KGS",
		},
		"limit": 3,
	}

	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var out struct {
		FilterURL string `json:"filter_url"`
		Top3      []any  `json:"top3"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)

	reply := "Вот, что удалось найти по вашему запросу. Хотите открыть подборку?"
	extra := map[string]any{
		"filter_url": out.FilterURL,
		"top3":       out.Top3,
	}
	return reply, extra, nil
}

// handleSeller — вызывает seller webhook (n8n).
func (s *service) handleSeller(ctx context.Context, sessionID, text string) (string, map[string]any, error) {
	url := s.cfg.N8NSellerWebhookURL
	if url == "" {
		return "Seller webhook не настроен", nil, nil
	}

	payload := map[string]any{
		"description": text,
	}

	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)

	reply := "Подготовил черновик объявления. Перейдите на страницу, чтобы опубликовать."
	return reply, out, nil
}

func extractLastUserText(_ *http.Request) string {
	// Пока у тебя нет истории в БД — можно возвращать последний текст из теста.
	return "Хочу телефон до 15000 сомов"
}
