package main

import (
	"context"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	// 👇 проверь пути: они должны совпадать со структурой твоего проекта
	"github.com/btynybekov/marketplace/config"      // если у тебя internal/config — замени импорт
	"github.com/btynybekov/marketplace/internal/ai" // <-- ВАЖНО: ai в internal/ai
	"github.com/btynybekov/marketplace/internal/handlers/factory"
	"github.com/btynybekov/marketplace/internal/repository"
	"github.com/btynybekov/marketplace/storage" // если у тебя internal/storage — замени импорт
)

func main() {
	// 1) Конфиг
	cfg := config.Load()

	// 2) БД
	ctx := context.Background()
	db, err := storage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer db.Close()

	// 3) Репозитории
	repos := repository.New(db.Pool)

	// 4) AI-клиент (openai | local)
	var aiClient ai.Client
	switch cfg.AIProvider {
	case "local":
		headers := map[string]string{}
		if cfg.LocalAIKey != "" {
			headers["X-API-Key"] = cfg.LocalAIKey
		}
		aiClient = ai.NewLocal(cfg.LocalAIURL, headers)
	default:
		aiClient = ai.NewOpenAI(cfg.OpenAIKey)
	}

	// 5) Шаблоны (если нужны HTML-страницы)
	var tmpl *template.Template
	// tmpl = factory.MustParseTemplates("web/templates/*.html")

	// 6) Фабрика хендлеров (сюда ПЕРЕДАЁМ aiClient)
	hf := factory.NewHandlersFactory(repos, tmpl, cfg, aiClient)

	// 7) Роутер и регистрация маршрутов
	r := mux.NewRouter()

	// Если нужен API-префикс и/или внешние middleware — раскомментируй:
	// api := r.PathPrefix("/api/v1").Subrouter()
	// api.Use(AuthMiddleware)
	// hf.RegisterRoutes(api)
	// А если префикс не нужен:
	hf.RegisterRoutes(r)

	// 8) HTTP-сервер + graceful shutdown
	srv := &http.Server{
		Addr:              ":" + cfg.PORT,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	go func() {
		log.Printf("Server listening on %s (env=%s)", srv.Addr, cfg.AppEnv)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve error: %v", err)
		}
	}()

	// ожидание сигнала для мягкой остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}
