package main

import (
	"log"
	"net/http"

	"github.com/btynybekov/marketplace/internal/repository"
	"github.com/btynybekov/marketplace/internal/services"

	"github.com/btynybekov/marketplace/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := repository.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Сервисы
	userService := services.NewUserService(userRepo)
	itemService := services.NewItemService(itemRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	itemHandler := handlers.NewItemHandler(itemService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	r := mux.NewRouter()

	// User endpoints
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Item endpoints
	r.HandleFunc("/items", itemHandler.List).Methods("GET")
	r.HandleFunc("/items/{id}", itemHandler.Get).Methods("GET")

	// Category endpoints
	r.HandleFunc("/categories", categoryHandler.List).Methods("GET")
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
