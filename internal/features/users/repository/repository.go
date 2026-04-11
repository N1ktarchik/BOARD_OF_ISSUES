package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pgUniqueViolation = "23505"
)

type UsersRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewUsersRepository(pool *pgxpool.Pool, log *slog.Logger) *UsersRepository {
	return &UsersRepository{
		pool: pool,
		log:  log,
	}
}
