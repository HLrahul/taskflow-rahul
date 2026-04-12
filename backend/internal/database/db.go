package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitDB() *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("Cannot ping database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully connected to Postgres!")
	runMigrations(connStr)

	return pool
}

func runMigrations(connStr string) {
	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		slog.Error("Could not create migrate instance", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Could not run up migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("Database migrations ran successfully (or no new changes)")
}
