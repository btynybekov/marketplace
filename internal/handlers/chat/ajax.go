package chat

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/btynybekov/marketplace/internal/repository"
)

type ChatAjaxHandler struct {
	Repo repository.Repository
}

func NewChatAjaxHandler(repo repository.Repository) *ChatAjaxHandler {
	return &ChatAjaxHandler{Repo: repo}
}
func (h *ChatAjaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		req.UserID = "anonymous"
	}

	ctx := r.Context()

	if err := h.Repo.SaveMessage(ctx, req.UserID, "user", req.Message); err != nil {
		http.Error(w, "failed to save message", http.StatusInternalServerError)
		log.Printf("SaveMessage error: %v", err)
		return
	}

	n8nBody := map[string]string{
		"userId":  req.UserID,
		"message": req.Message,
	}
	bodyBytes, _ := json.Marshal(n8nBody)

	resp, err := http.Post("http://localhost:5678/webhook-test/chat", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		http.Error(w, "failed to call n8n webhook", http.StatusInternalServerError)
		log.Printf("n8n webhook error: %v", err)
		return
	}
	defer resp.Body.Close()

	var n8nResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&n8nResp); err != nil {
		http.Error(w, "invalid response from n8n", http.StatusInternalServerError)
		log.Printf("n8n response decode error: %v", err)
		return
	}

	aiMessage, ok := n8nResp["message"].(string)
	if !ok {
		aiMessage = "Заглушка AI: сообщение отсутствует"
	}

	// Сохраняем ответ AI в БД
	if err := h.Repo.SaveMessage(ctx, req.UserID, "assistant", aiMessage); err != nil {
		http.Error(w, "failed to save AI message", http.StatusInternalServerError)
		log.Printf("SaveMessage AI error: %v", err)
		return
	}

	// Возвращаем ответ пользователю
	respJSON := models.ChatResponse{
		UserID:    req.UserID,
		Message:   aiMessage,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respJSON)
}
