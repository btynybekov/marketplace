package ai

import "context"

// Message — единый формат сообщений для любых LLM.
type Message struct {
	Role    string // "system" | "user" | "assistant" | "tool"
	Content string
}

// Client — абстракция поверх любого поставщика LLM.
type Client interface {
	Chat(ctx context.Context, model string, temperature float64, messages []Message) (string, error)
}
