package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"taskflow-backend/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, name, email, hashedPassword string) (*models.User, error) {
	var user models.User

	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at
	`

	err := r.db.QueryRow(ctx, query, name, email, hashedPassword).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
