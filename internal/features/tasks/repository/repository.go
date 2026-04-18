package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TasksRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewTasksRepository(pool *pgxpool.Pool, log *slog.Logger) *TasksRepository {
	return &TasksRepository{
		pool: pool,
		log:  log,
	}
}
