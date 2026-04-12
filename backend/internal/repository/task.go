package repository

import (
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"taskflow-backend/internal/models"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.ProjectID, t.AssigneeID, t.DueDate).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)

	return t, err
}

func (r *TaskRepository) GetTasks(ctx context.Context, projectID, status, assigneeID string, limit, offset int) ([]models.Task, error) {
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE project_id = $1`

	args := []interface{}{projectID}
	argCount := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if assigneeID != "" {
		query += fmt.Sprintf(" AND assignee_id = $%d", argCount)
		args = append(args, assigneeID)
		argCount++
	}
	query += " ORDER BY created_at DESC LIMIT $%d OFFSET $%d"
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.AssigneeID, t.DueDate, t.ID).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	return t, err
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}

func (r *TaskRepository) CheckProjectOwnerByTask(ctx context.Context, taskID string, userID string) (bool, error) {
	var ownerID string
	query := `
		SELECT p.owner_id 
		FROM projects p 
		JOIN tasks t ON p.id = t.project_id 
		WHERE t.id = $1
	`

	err := r.db.QueryRow(ctx, query, taskID).Scan(&ownerID)
	if err != nil {
		return false, err
	}

	return ownerID == userID, nil
}
