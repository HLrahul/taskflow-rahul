package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"taskflow-backend/internal/database"
	"taskflow-backend/internal/handler"
	taskflowMiddleware "taskflow-backend/internal/middleware"
	"taskflow-backend/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load("../.env")
	if err != nil {
		slog.Warn("No .env file found or couldn't load it. Relying on system environment variables.")
	}

	dbPool := database.InitDB()
	defer dbPool.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health Check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "TaskFlow API is running!"})
	})

	// Auth Routes
	userRepo := repository.NewUserRepository(dbPool)
	authHandler := handler.NewAuthHandler(userRepo)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	projectRepo := repository.NewProjectRepository(dbPool)
	taskRepo := repository.NewTaskRepository(dbPool)

	// Project Routes
	projectHandler := handler.NewProjectHandler(projectRepo, taskRepo)

	r.Group(func(r chi.Router) {
		r.Use(taskflowMiddleware.RequireAuth)
		r.Get("/projects", projectHandler.ListProjects)
		r.Post("/projects", projectHandler.CreateProject)
		r.Get("/projects/{id}", projectHandler.GetProject)
		r.Patch("/projects/{id}", projectHandler.UpdateProject)
		r.Delete("/projects/{id}", projectHandler.DeleteProject)
		r.Get("/projects/{id}/stats", projectHandler.GetProjectStats)
	})

	// Task Routes
	taskHandler := handler.NewTaskHandler(taskRepo, projectRepo)

	r.Group(func(r chi.Router) {
		r.Use(taskflowMiddleware.RequireAuth)
		r.Get("/projects/{id}/tasks", taskHandler.ListTasks)
		r.Post("/projects/{id}/tasks", taskHandler.CreateTask)
		r.Patch("/tasks/{id}", taskHandler.UpdateTask)
		r.Delete("/tasks/{id}", taskHandler.DeleteTask)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func () {
		slog.Info("Starting API server", "port", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("Received SIGTERM/SIGINT. Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server successfully stopped.")
}
