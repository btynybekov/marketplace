package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type OpenAIClient struct {
	Key    string
	Client *http.Client
}

func NewOpenAI(key string) *OpenAIClient {
	if key == "" {
		panic("OPENAI_API_KEY is empty")
	}
	return &OpenAIClient{
		Key:    key,
		Client: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *OpenAIClient) Chat(ctx context.Context, model string, temperature float64, messages []Message) (string, error) {
	payload := map[string]any{
		"model":       model,
		"temperature": temperature,
		"messages":    messages,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", errors.New("openai http status: " + resp.Status)
	}

	var out struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("openai: empty choices")
	}
	return out.Choices[0].Message.Content, nil
}
