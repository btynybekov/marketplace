package handlers

import (
	"encoding/json"
	"marketplace/internal/models"
	"marketplace/internal/services"
	"net/http"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	json.NewDecoder(r.Body).Decode(&cat)
	if err := h.service.Create(&cat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, _ := h.service.List()
	json.NewEncoder(w).Encode(cats)
}
