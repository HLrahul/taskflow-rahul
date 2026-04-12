package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"taskflow-backend/internal/models"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, name, description, ownerID string) (*models.Project, error) {
	var p models.Project

	query := `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, owner_id , created_at
	`

	err := r.db.QueryRow(ctx, query, name, description, ownerID).Scan(
		&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt,
	)

	return &p, err
}

func (r *ProjectRepository) GetUserProjects(ctx context.Context, userId string, limit, offset int) ([]models.Project, error) {
	query := `
		SELECT p.id, p.name, p.description, p.owner_id, p.created_at
		FROM projects p
		WHERE p.owner_id = $1 OR EXISTS (
			SELECT 1 FROM tasks t WHERE t.project_id = p.id AND t.assignee_id = $1
		)
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	projects := make([]models.Project, 0)

	for rows.Next() {
		var p models.Project

		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *ProjectRepository) CountUserProjects(ctx context.Context, userId string) (int, error) {
	query := `
		SELECT COUNT(*) FROM projects p
		WHERE p.owner_id = $1 OR EXISTS (
			SELECT 1 FROM tasks t WHERE t.project_id = p.id AND t.assignee_id = $1
		)
	`
	var count int
	err := r.db.QueryRow(ctx, query, userId).Scan(&count)
	return count, err
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id, ownerID string) (*models.Project, error) {
	var p models.Project

	query := `
		SELECT p.id, p.name, p.description, p.owner_id, p.created_at
		FROM projects p
		WHERE p.id = $1 AND (p.owner_id = $2 OR EXISTS (
			SELECT 1 FROM tasks t WHERE t.project_id = p.id AND t.assignee_id = $2
		))
	`

	err := r.db.QueryRow(ctx, query, id, ownerID).Scan(
		&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, id, name, description, ownerID string) (*models.Project, error) {
	var p models.Project

	query := `
		UPDATE projects 
		SET name = $1, description = $2 
		WHERE id = $3 AND owner_id = $4
		RETURNING id, name, description, owner_id, created_at
	`

	err := r.db.QueryRow(ctx, query, name, description, id, ownerID).Scan(
		&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id, ownerID string) error {
	query := `
		DELETE FROM projects
		WHERE id = $1 AND owner_id = $2
	`

	tag, err := r.db.Exec(ctx, query, id, ownerID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return context.DeadlineExceeded
	}

	return nil
}

func (r *ProjectRepository) GetProjectStats(ctx context.Context, projectID string) (*models.ProjectStats, error) {
	stats := &models.ProjectStats{
		ByStatus:   make(map[string]int),
		ByAssignee: make(map[string]int),
	}

	statusQuery := `
		SELECT status, COUNT(*) FROM tasks
		WHERE project_id = $1 GROUP BY status
	`

	statusRows, err := r.db.Query(ctx, statusQuery, projectID)
	if err == nil {
		defer statusRows.Close()

		for statusRows.Next() {
			var status string
			var count int
			statusRows.Scan(&status, &count)
			stats.ByStatus[status] = count
		}
	}

	assigneeQuery := `
		SELECT COALESCE(assignee_id::text, 'unassigned'), COUNT(*) 
		FROM tasks 
		WHERE project_id = $1 
		GROUP BY assignee_id
	`

	assigneeRows, err := r.db.Query(ctx, assigneeQuery, projectID)
	if err == nil {
		defer assigneeRows.Close()

		for assigneeRows.Next() {
			var assignee string
			var count int
			assigneeRows.Scan(&assignee, &count)
			stats.ByAssignee[assignee] = count
		}
	}

	return stats, nil
}
