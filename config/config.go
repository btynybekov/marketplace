package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// EnvConfig содержит все переменные окружения приложения.
type EnvConfig struct {
	// Системные
	AppEnv string // "development" | "production"
	PORT   string // порт HTTP-сервера

	// Подключения
	DatabaseURL string // postgres://user:pass@host:port/db?sslmode=disable

	// Сервисы
	AIProvider    string // "openai" | "local"
	AIModel       string
	AITemperature float64
	OpenAIKey     string
	LocalAIURL    string // http://ollama:11434 (пример)
	LocalAIKey    string // если нужно

	N8NBuyerWebhookURL  string // webhook ассистента покупателя
	N8NSellerWebhookURL string // webhook ассистента продавца
	AssetsBaseURL       string // базовый URL для статики или CDN

	// Настройки (опционально)
	LogLevel  string // info | debug | warn
	DebugMode bool   // включить подробные логи
}

// Env — глобальная конфигурация, доступная из любой точки
var Env EnvConfig

// Load инициализирует конфигурацию из .env и системных переменных
func Load() EnvConfig {
	// Загружаем .env (если есть)
	_ = godotenv.Load() // не падаем, если файла нет

	cfg := EnvConfig{
		AppEnv:              getenvOrDefault("APP_ENV", "development"),
		PORT:                getenvOrDefault("PORT", "8080"),
		DatabaseURL:         getenvOrDefault("DATABASE_URL", ""),
		AIProvider:          getenvOrDefault("AI_PROVIDER", "openai"),
		AIModel:             getenvOrDefault("AI_MODEL", "gpt-4o-mini"),
		AITemperature:       getenvAsFloat("AI_TEMPERATURE", 0.2),
		OpenAIKey:           getenvOrDefault("OPENAI_API_KEY", ""),
		LocalAIURL:          getenvOrDefault("LOCAL_AI_URL", ""),
		LocalAIKey:          getenvOrDefault("LOCAL_AI_KEY", ""),
		N8NBuyerWebhookURL:  getenvOrDefault("N8N_BUYER_ASSISTANT_WEBHOOK_URL", ""),
		N8NSellerWebhookURL: getenvOrDefault("N8N_SELLER_ASSISTANT_WEBHOOK_URL", ""),
		AssetsBaseURL:       getenvOrDefault("ASSETS_BASE_URL", "/static"),
		LogLevel:            getenvOrDefault("LOG_LEVEL", "info"),
		DebugMode:           getenvAsBool("DEBUG", false),
	}

	// Валидация обязательных параметров
	if cfg.DatabaseURL == "" {
		log.Fatal("[CONFIG] DATABASE_URL is not set — cannot connect to DB")
	}

	Env = cfg
	return cfg
}

// GetDSN возвращает строку подключения к БД
func GetDSN() string {
	if Env.DatabaseURL == "" {
		log.Fatal("[CONFIG] DATABASE_URL is empty")
	}
	return Env.DatabaseURL
}

//
// ──────────────────────────────────────────────
//   ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ──────────────────────────────────────────────
//

// getenvOrDefault возвращает значение переменной окружения или дефолт.
func getenvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// getenvAsBool — парсит bool ("true"/"1" → true).
func getenvAsBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}

func getenvAsFloat(key string, defaultVal float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return defaultVal
	}
	return float64(f)
}
