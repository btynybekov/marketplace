package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/btynybekov/marketplace/internal/services"

	"github.com/gorilla/mux"
)

type ItemHandler struct {
	service *services.ItemService
}

func NewItemHandler(service *services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	json.NewDecoder(r.Body).Decode(&item)
	userID := r.Context().Value("userID").(int64)
	item.UserID = userID
	if err := h.service.Create(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	item, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	json.NewDecoder(r.Body).Decode(&item)
	item.ID, _ = strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	userID := r.Context().Value("userID").(int64)
	if err := h.service.Update(&item, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	userID := r.Context().Value("userID").(int64)
	if err := h.service.Delete(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, _ := h.service.List()
	json.NewEncoder(w).Encode(items)
}
