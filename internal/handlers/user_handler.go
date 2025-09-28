package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/btynybekov/marketplace/internal/models"
	"github.com/btynybekov/marketplace/internal/services"
	"github.com/btynybekov/marketplace/pkg/auth"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register endpoint
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := h.service.Register(user, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, _ := auth.GenerateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Login endpoint
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Profile endpoint (protected)
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	user, err := h.service.GetByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
