package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/middleware"
	"taskflow-backend/internal/models"
	"taskflow-backend/internal/repository"
)

type ProjectHandler struct {
	repo     *repository.ProjectRepository
	taskRepo *repository.TaskRepository
}

func NewProjectHandler(repo *repository.ProjectRepository, taskRepo *repository.TaskRepository) *ProjectHandler {
	return &ProjectHandler{repo: repo, taskRepo: taskRepo}
}

// GET /projects
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)
	limit, offset := parsePagination(r)

	projects, err := h.repo.GetUserProjects(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	total, _ := h.repo.CountUserProjects(r.Context(), userID)

	page := (offset / limit) + 1

	json.NewEncoder(w).Encode(map[string]interface{}{
		"projects": projects,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// POST /projects
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid json payload"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "validation failed", "fields": {"name": "is required"}}`))
		return
	}

	project, err := h.repo.CreateProject(r.Context(), req.Name, req.Description, userID)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// GET /projects/:id
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)

	projectID := chi.URLParam(r, "id")
	project, err := h.repo.GetProjectByID(r.Context(), projectID, userID)
	if err != nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	limit, offset := parsePagination(r)
	tasks, err := h.taskRepo.GetTasks(r.Context(), projectID, "", "", limit, offset)
	if err != nil {
		slog.Error("GetTasks failed", "error", err, "projectID", projectID)
		tasks = []models.Task{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          project.ID,
		"name":        project.Name,
		"description": project.Description,
		"owner_id":    project.OwnerID,
		"tasks":       tasks,
	})
}

// PATCH /projects/:id
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)
	projectID := chi.URLParam(r, "id")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid json payload"}`, http.StatusBadRequest)
		return
	}

	project, err := h.repo.UpdateProject(r.Context(), projectID, req.Name, req.Description, userID)
	if err != nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(project)
}

// DELETE /projects/:id
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	projectID := chi.URLParam(r, "id")

	err := h.repo.DeleteProject(r.Context(), projectID, userID)
	if err != nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /projects/:id/stats
func (h *ProjectHandler) GetProjectStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)
	projectID := chi.URLParam(r, "id")

	_, err := h.repo.GetProjectByID(r.Context(), projectID, userID)
	if err != nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	stats, err := h.repo.GetProjectStats(r.Context(), projectID)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch stats"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func parsePagination(r *http.Request) (int, int) {
	limit := 10
	page := 1

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	offset := (page - 1) * limit
	return limit, offset
}
