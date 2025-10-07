package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/btynybekov/marketplace/config"
	db "github.com/btynybekov/marketplace/database"
	"github.com/btynybekov/marketplace/internal/handlers/factory"
	"github.com/btynybekov/marketplace/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	// 1️⃣ Загружаем конфиг
	conf := config.LoadFromEnv()
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))
	// 2️⃣ Инициализируем базу данных
	db.InitDB(conf)
	defer db.CloseDB()

	// 3️⃣ Создаём репозиторий
	repo := repository.NewPostgresRepo(db.Pool)

	// 4️⃣ Парсим шаблоны один раз
	tmpl := factory.MustParseTemplates("templates/*.html")

	// 5️⃣ Создаём фабрику хендлеров
	factory := factory.NewHandlersFactory(repo, tmpl, conf)

	// 6️⃣ Создаём роутер Gorilla Mux
	r := mux.NewRouter()

	// 7️⃣ Регистрируем маршруты через фабрику
	factory.RegisterRoutes(r)

	// 8️⃣ Берём порт из env или используем 8080
	port := config.Env.PORT
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)

	// 9️⃣ Запускаем HTTP сервер
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
