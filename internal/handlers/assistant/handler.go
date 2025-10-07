package assistant

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type AssistantRequest struct {
	UserID string                 `json:"userId"`
	Task   string                 `json:"task"`
	Data   map[string]interface{} `json:"data"`
}

type AssistantResponse struct {
	Message string `json:"message"`
}

type AssistantHandler struct {
	N8nWebhookURL string
}

func NewAssistantHandler(url string) *AssistantHandler {
	return &AssistantHandler{N8nWebhookURL: url}
}

func (h *AssistantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req AssistantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		req.UserID = "anonymous"
	}

	bodyBytes, _ := json.Marshal(req)
	resp, err := http.Post(h.N8nWebhookURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Ошибка вызова n8n: %v", err)
		http.Error(w, "failed to call n8n", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var assistantResp AssistantResponse
	if err := json.NewDecoder(resp.Body).Decode(&assistantResp); err != nil {
		log.Printf("Ошибка декодирования ответа n8n: %v", err)
		http.Error(w, "invalid response from n8n", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assistantResp)
}
