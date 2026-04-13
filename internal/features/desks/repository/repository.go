package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pgForeignKeyViolation = "23503"
)

type DesksRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewDesksRepository(pool *pgxpool.Pool, log *slog.Logger) *DesksRepository {
	return &DesksRepository{
		pool: pool,
		log:  log,
	}
}
