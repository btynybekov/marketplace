package chat

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/btynybekov/marketplace/internal/handlers/shared"
)

//
// ─── HTTP DTO ───────────────────────────────────────────────────────────────────
//

type startSessionReq struct {
	SessionID string `json:"session_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}
type startSessionResp struct {
	SessionID string `json:"session_id"`
}

type sendMessageReq struct {
	SessionID string            `json:"session_id"`
	Text      string            `json:"text"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// messageDTO используется и в service.go (тот же пакет chat)
type messageDTO struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"` // "user" | "assistant"
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type sendMessageResp struct {
	Reply       messageDTO   `json:"reply"`
	Top3        []any        `json:"top3,omitempty"`       // если был поиск (buyer)
	FilterURL   string       `json:"filter_url,omitempty"` // если был поиск (buyer)
	NewMessages []messageDTO `json:"new_messages,omitempty"`
}

type historyResp struct {
	SessionID string       `json:"session_id"`
	Messages  []messageDTO `json:"messages"`
}

//
// ─── HTTP HANDLER ───────────────────────────────────────────────────────────────
//

// ChatHandler — HTTP-обёртка над сервисом чата (бизнес-логика в service.go).
type ChatHandler struct {
	svc Service
}

// NewChatHTTP — конструктор HTTP-хендлера для API чата.
func NewChatHTTP(svc Service) *ChatHandler { return &ChatHandler{svc: svc} }

// StartSession — POST /chat/session
func (h *ChatHandler) StartSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req startSessionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.BadRequest(w, "invalid JSON")
			return
		}

		sid, err := h.svc.StartSession(r, req.UserID, req.SessionID)
		if err != nil {
			shared.InternalError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, startSessionResp{SessionID: sid})
	})
}

// SendMessage — POST /chat/ajax
// Принимает текст пользователя, сохраняет его (если реализовано в сервисе),
// генерирует ответ ассистента: обычная беседа или вызов buyer/seller по контексту.
func (h *ChatHandler) SendMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req sendMessageReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.BadRequest(w, "invalid JSON")
			return
		}
		if req.SessionID == "" || req.Text == "" {
			shared.BadRequest(w, "session_id and text are required")
			return
		}

		// сохраняем сообщение пользователя (опционально — внутри сервиса)
		if _, err := h.svc.AppendUserMessage(r, req.SessionID, req.Text, req.Meta); err != nil {
			shared.InternalError(w, err)
			return
		}

		// генерируем ответ (LLM или buyer/seller ветка)
		reply, extra, err := h.svc.GenerateAssistantReply(r, req.SessionID)
		if err != nil {
			shared.InternalError(w, err)
			return
		}

		resp := sendMessageResp{
			Reply: messageDTO{
				ID:        "", // если сервис вернёт ID — подставь сюда
				Role:      "assistant",
				Text:      reply,
				CreatedAt: time.Now().UTC(),
			},
		}

		// пробрасываем дополнительные поля от сервиса (например, топ-3 и filter_url)
		if extra != nil {
			if v, ok := extra["top3"]; ok {
				if top, ok2 := v.([]any); ok2 {
					resp.Top3 = top
				}
			}
			if v, ok := extra["filter_url"]; ok {
				if s, ok2 := v.(string); ok2 {
					resp.FilterURL = s
				}
			}
			// при желании — пробрось и другие поля из extra...
		}

		shared.WriteJSON(w, http.StatusOK, resp)
	})
}

// GetHistory — GET /chat/history или /chat/history/{session_id}
func (h *ChatHandler) GetHistory() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// поддерживаем и query, и path-переменную
		sid := r.URL.Query().Get("session_id")
		if sid == "" {
			if v, ok := mux.Vars(r)["session_id"]; ok {
				sid = v
			}
		}
		if sid == "" {
			shared.BadRequest(w, "session_id is required")
			return
		}

		msgs, err := h.svc.GetHistory(r, sid, 50)
		if err != nil {
			shared.InternalError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, historyResp{
			SessionID: sid,
			Messages:  msgs,
		})
	})
}
