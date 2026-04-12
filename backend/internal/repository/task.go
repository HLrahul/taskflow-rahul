package repository

import (
	"context"
	"fmt"

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
		WITH inserted_task AS (
			INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		)
		SELECT it.*, u.name, u.email 
		FROM inserted_task it 
		LEFT JOIN users u ON it.assignee_id = u.id
	`

	err := r.db.QueryRow(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.ProjectID, t.AssigneeID, t.DueDate).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt, &t.AssigneeName, &t.AssigneeEmail,
	)

	return t, err
}

func (r *TaskRepository) GetTasks(ctx context.Context, projectID, status, assigneeID string, limit, offset int) ([]models.Task, error) {
	query := `
		SELECT t.id, t.title, t.description, t.status, t.priority, t.project_id, t.assignee_id, 
		       u.name, u.email, t.due_date, t.created_at, t.updated_at 
		FROM tasks t 
		LEFT JOIN users u ON t.assignee_id = u.id 
		WHERE t.project_id = $1
	`

	args := []interface{}{projectID}
	argCount := 2

	if status != "" {
		query += fmt.Sprintf(" AND t.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if assigneeID != "" {
		query += fmt.Sprintf(" AND t.assignee_id = $%d", argCount)
		args = append(args, assigneeID)
		argCount++
	}
	query += fmt.Sprintf(" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.AssigneeName, &t.AssigneeEmail, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	query := `
		WITH updated_task AS (
			UPDATE tasks 
			SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
			RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		)
		SELECT ut.*, u.name, u.email 
		FROM updated_task ut 
		LEFT JOIN users u ON ut.assignee_id = u.id
	`

	err := r.db.QueryRow(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.AssigneeID, t.DueDate, t.ID).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt, &t.AssigneeName, &t.AssigneeEmail,
	)
	return t, err
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}

// GetTaskAccessLevel returns (isOwner, isAssignee, err)
func (r *TaskRepository) GetTaskAccessLevel(ctx context.Context, taskID string, userID string) (bool, bool, error) {
	var ownerID string
	var assigneeID *string

	query := `
		SELECT p.owner_id, t.assignee_id 
		FROM projects p 
		JOIN tasks t ON p.id = t.project_id 
		WHERE t.id = $1
	`

	err := r.db.QueryRow(ctx, query, taskID).Scan(&ownerID, &assigneeID)
	if err != nil {
		return false, false, err
	}

	isOwner := ownerID == userID
	isAssignee := assigneeID != nil && *assigneeID == userID

	return isOwner, isAssignee, nil
}

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, taskID string, status string) (*models.Task, error) {
	query := `
		WITH updated_status AS (
			UPDATE tasks 
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
			RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		)
		SELECT us.*, u.name, u.email 
		FROM updated_status us 
		LEFT JOIN users u ON us.assignee_id = u.id
	`
	var t models.Task
	err := r.db.QueryRow(ctx, query, status, taskID).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt, &t.AssigneeName, &t.AssigneeEmail,
	)
	return &t, err
}
