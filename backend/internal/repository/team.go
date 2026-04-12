package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"taskflow-backend/internal/models"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) AddTeamMember(ctx context.Context, ownerID string, memberID string) error {
	query := `
		INSERT INTO team_members (owner_id, member_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, ownerID, memberID)
	return err
}

func (r *TeamRepository) GetTeamMembers(ctx context.Context, ownerID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.created_at
		FROM users u
		JOIN team_members tm ON u.id = tm.member_id
		WHERE tm.owner_id = $1
	`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	return members, nil
}
