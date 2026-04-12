package handler

import (
	"net/http"
	"encoding/json"

	"github.com/go-chi/chi/v5"

	"taskflow-backend/internal/models"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/internal/repository"
)

type TaskHandler struct {
	repo        *repository.TaskRepository
	projectRepo *repository.ProjectRepository
}

func NewTaskHandler(repo *repository.TaskRepository, projectRepo *repository.ProjectRepository) *TaskHandler {
	return &TaskHandler{repo: repo, projectRepo: projectRepo}
}

// GET /projects/:id/tasks?status=todo&assignee=uuid
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID := chi.URLParam(r, "id")
	status := r.URL.Query().Get("status")
	assigneeID := r.URL.Query().Get("assignee")

	limit, offset := parsePagination(r)

	tasks, err := h.repo.GetTasks(r.Context(), projectID, status, assigneeID, limit, offset)
	if err != nil {
		http.Error(w, `{"error": "Failed to list tasks"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

// POST /projects/:id/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)
	projectID := chi.URLParam(r, "id")

	_, err := h.projectRepo.GetProjectByID(r.Context(), projectID, userID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Forbidden"}`))
		return
	}

	var req models.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid json payload"}`, http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Validation failed", "fields": {"title": "is required"}}`))
		return
	}

	if req.Status == "" {
		req.Status = "todo"
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}
	req.ProjectID = projectID

	task, err := h.repo.CreateTask(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error": "Failed to create task"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// PATCH /tasks/:id
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIDKey).(string)
	taskID := chi.URLParam(r, "id")

	isOwner, isAssignee, err := h.repo.GetTaskAccessLevel(r.Context(), taskID, userID)
	if err != nil || (!isOwner && !isAssignee) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Forbidden"}`))
		return
	}

	var req models.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid json payload"}`, http.StatusBadRequest)
		return
	}
	req.ID = taskID

	var task *models.Task
	if !isOwner && isAssignee {
		// Strict lock: Assignees can only patch the status.
		task, err = h.repo.UpdateTaskStatus(r.Context(), taskID, req.Status)
	} else {
		// Owner has full modification rights
		task, err = h.repo.UpdateTask(r.Context(), &req)
	}

	if err != nil {
		http.Error(w, `{"error": "Failed to update task"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	taskID := chi.URLParam(r, "id")

	isOwner, _, err := h.repo.GetTaskAccessLevel(r.Context(), taskID, userID)
	if err != nil || !isOwner {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Forbidden"}`))
		return
	}

	err = h.repo.DeleteTask(r.Context(), taskID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
