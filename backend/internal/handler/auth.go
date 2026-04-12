package handler

import (
	"net/http"
	"encoding/json"

	"taskflow-backend/pkg/utils"
	"taskflow-backend/internal/models"
	"taskflow-backend/internal/repository"
)

type AuthHandler struct {
	repo *repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid json payload"}`, http.StatusBadRequest)
		return
	}

	fields := make(map[string]string)

	if req.Name == "" {
		fields["name"] = "is required"
	}

	if req.Email == "" {
		fields["email"] = "is required"
	}

	if req.Password == "" {
		fields["password"] = "is required"
	}

	if len(fields) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "validation failed",
			"fields": fields,
		})

		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	user, err := h.repo.CreateUser(r.Context(), req.Name, req.Email, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "validation failed", "fields": {"email": "already exists"}}`))
		return
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid json payload"}`, http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "No user with this email"}`))
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Incorrect password"}`))
		return
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}
