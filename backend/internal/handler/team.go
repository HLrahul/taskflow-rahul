package handler

import (
	"encoding/json"
	"net/http"

	"taskflow-backend/internal/middleware"
	"taskflow-backend/internal/repository"
)

type TeamHandler struct {
	teamRepo *repository.TeamRepository
	userRepo *repository.UserRepository
}

func NewTeamHandler(teamRepo *repository.TeamRepository, userRepo *repository.UserRepository) *TeamHandler {
	return &TeamHandler{teamRepo: teamRepo, userRepo: userRepo}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value(middleware.UserIDKey).(string)

	members, err := h.teamRepo.GetTeamMembers(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch team"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"members": members})
}

func (h *TeamHandler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid json payload"}`, http.StatusBadRequest)
		return
	}

	userToAdd, err := h.userRepo.GetUserByEmail(r.Context(), req.Email)
	if err != nil || userToAdd == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "user with this email not found"}`))
		return
	}

	if userToAdd.ID == userID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "you cannot add yourself to your team"}`))
		return
	}

	err = h.teamRepo.AddTeamMember(r.Context(), userID, userToAdd.ID)
	if err != nil {
		http.Error(w, `{"error": "failed to add team member"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "team member added successfully!", "member": userToAdd})
}
