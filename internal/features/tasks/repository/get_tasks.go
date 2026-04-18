package repository

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *TasksRepository) GetTasksFromOneDesk(ctx context.Context, filter *domain.TaskFilter) ([]*domain.Task, error) {
	r.log.Info("getting tasks from database")

	query := `
		SELECT t.id, t.author_id, t.desk_id, t.name, t.description, t.done, t.deadline, t.created_at
		FROM tasks t
		JOIN desk_members dm ON t.desk_id = dm.desk_id
		WHERE dm.user_id = $1 AND t.desk_id = $2`

	args := []any{filter.UserId, filter.DeskId}
	argID := 3

	if filter.Done != nil {
		query += fmt.Sprintf(" AND t.done = $%d", argID)
		args = append(args, *filter.Done)
		argID++
	}

	query += " ORDER BY t.created_at DESC"

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("err", err))
		return nil, core_errors.ServerError()
	}
	defer rows.Close()

	tasks := make([]*domain.Task, 0, filter.Limit)

	for rows.Next() {
		var t taskModel
		if err := t.scan(rows); err != nil {
			r.log.Error("row scan failed", slog.Any("err", err))
			return nil, core_errors.ServerError()
		}
		domainTask := modelToDomain(t)
		tasks = append(tasks, &domainTask)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("rows cursor error", slog.Any("err", err))
		return nil, core_errors.ServerError()
	}

	r.log.Info("tasks retrieved", slog.Int("count", len(tasks)))
	return tasks, nil
}

func (r *TasksRepository) GetTaskByID(ctx context.Context, taskID, userID uuid.UUID) (*domain.Task, error) {
	r.log.Info("get task by id from database")

	query := `SELECT t.id, t.author_id, t.desk_id, t.name, t.description, t.done, t.deadline, t.created_at 
			FROM tasks t
			JOIN desk_members dm ON t.desk_id = dm.desk_id
			WHERE dm.user_id = $1 AND t.id = $2`

	row := r.pool.QueryRow(ctx, query, userID, taskID)

	task := taskModel{}

	if err := task.scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn("user have not access to desk or task not faund", slog.Any("userID", userID), slog.Any("taskID", taskID))

			return nil, core_errors.TaskNotFound()
		}

		r.log.Error("repository error", slog.Any("taskID", taskID), slog.Any("userID", userID), slog.Any("err", err))

		return nil, core_errors.ServerError()
	}

	domainTask := modelToDomain(task)

	r.log.Info("task retrieved", slog.Any("taskID", taskID))
	return &domainTask, nil

}
