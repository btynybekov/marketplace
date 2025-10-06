package main

import (
	"log"
	"net/http"

	"github.com/btynybekov/marketplace/config"
	db "github.com/btynybekov/marketplace/database"
	"github.com/btynybekov/marketplace/internal/handlers"
	"github.com/btynybekov/marketplace/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig("configs/config.yaml")
	db.InitDB()
	defer db.CloseDB()

	// Создаём репозиторий
	repo := repository.NewPostgresRepo(db.Pool)

	// Создаём handler через Repo
	h := handlers.NewHandler(repo)

	router := mux.NewRouter()
	router.HandleFunc("/", h.HomePage)
	router.HandleFunc("/category/{slug}", h.CategoryPage)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
