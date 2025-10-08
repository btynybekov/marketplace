package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// EnvConfig хранит все глобальные переменные окружения
type EnvConfig struct {
	PORT                string // Порт сервера
	DatabaseURL         string // Полная строка подключения к БД
	N8NBuyerWebhookURL  string // Вебхук ассистента покупателя
	N8NSellerWebhookURL string // Вебхук ассистента продавца
	AssetsBaseURL       string // Базовый URL для статических файлов
}

// Env — глобальная структура для всего приложения
var Env EnvConfig

// LoadFromEnv загружает переменные окружения и .env в EnvConfig
func LoadFromEnv() EnvConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found, reading from system environment")
	}
	log.Println(getenvOrDefault("DATABASE_URL", ""))

	cfg := EnvConfig{
		// Серверный порт
		PORT: getenvOrDefault("PORT", "8080"),
		// DATABASE_URL
		DatabaseURL:         getenvOrDefault("DATABASE_URL", ""),
		N8NBuyerWebhookURL:  getenvOrDefault("N8N_BUYER_ASSISTANT_WEBHOOK_URL", ""),
		N8NSellerWebhookURL: getenvOrDefault("N8N_SELLER_ASSISTANT_WEBHOOK_URL", ""),
		AssetsBaseURL:       getenvOrDefault("ASSETS_BASE_URL", "/static"),
	}
	return cfg
}

// GetDSN возвращает строку подключения к БД
func GetDSN(cfg EnvConfig) string {
	if cfg.DatabaseURL != "" {
		return cfg.DatabaseURL
	}
	log.Fatal("DATABASE_URL is not set, cannot connect to DB")
	return ""
}

// Вспомогательная функция для получения значения с дефолтом
func getenvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
