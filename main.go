package main

import (
	"log"
	"marketplace/internal/handlers"
	"marketplace/internal/repository"
	"marketplace/internal/services"
	"marketplace/pkg/auth"
	"net/http"

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
	r.Handle("/users/me", auth.JWTMiddleware(userHandler.Profile)).Methods("GET")

	// Item endpoints
	r.Handle("/items", auth.JWTMiddleware(itemHandler.Create)).Methods("POST")
	r.HandleFunc("/items", itemHandler.List).Methods("GET")
	r.HandleFunc("/items/{id}", itemHandler.Get).Methods("GET")
	r.Handle("/items/{id}", auth.JWTMiddleware(itemHandler.Update)).Methods("PUT")
	r.Handle("/items/{id}", auth.JWTMiddleware(itemHandler.Delete)).Methods("DELETE")

	// Category endpoints
	r.HandleFunc("/categories", categoryHandler.List).Methods("GET")
	r.Handle("/categories", auth.JWTMiddleware(categoryHandler.Create)).Methods("POST")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
