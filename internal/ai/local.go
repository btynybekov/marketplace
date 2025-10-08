package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// LocalClient — универсальный клиент для локального/самостоятельного HTTP API.
// Ожидается endpoint вида POST {BaseURL}/chat с payload:
// { "model": "...", "temperature": 0.2, "messages": [{role, content}, ...] }
// и ответ: { "reply": "..." }
type LocalClient struct {
	BaseURL string
	Client  *http.Client
	// Доп. заголовки/ключи если нужны (например, X-API-Key)
	Headers map[string]string
}

func NewLocal(baseURL string, headers map[string]string) *LocalClient {
	if baseURL == "" {
		panic("LocalAI BaseURL is empty")
	}
	return &LocalClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 20 * time.Second},
		Headers: headers,
	}
}

func (c *LocalClient) Chat(ctx context.Context, model string, temperature float64, messages []Message) (string, error) {
	payload := map[string]any{
		"model":       model,
		"temperature": temperature,
		"messages":    messages,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", errors.New("local llm http status: " + resp.Status)
	}

	var out struct {
		Reply string `json:"reply"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Reply, nil
}
