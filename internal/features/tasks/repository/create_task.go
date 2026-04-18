package repository

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (r *TasksRepository) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	r.log.Info("create task", slog.Any("authorID", task.AuthorId), slog.Any("deskID", task.DeskId))

	query := `INSERT INTO tasks (id,author_id,desk_id,name,description,done,deadline)
			SELECT $1, $2, $3, $4, $5, $6, $7
			WHERE EXISTS (
   	 			SELECT 1 FROM desk_members WHERE desk_id = $3 AND user_id = $2
			)
			RETURNING id,author_id,desk_id,name,description,done,deadline,created_at;`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Id,
		task.AuthorId,
		task.DeskId,
		task.Name,
		task.Description,
		task.Done,
		task.Deadline)

	saveTask := taskModel{}

	if err := saveTask.scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn("user have not access to desk", slog.Any("userID", task.AuthorId), slog.Any("deskID", task.DeskId))

			return nil, core_errors.UserHaveNotAccessToDesk(task.AuthorId.String(), task.DeskId.String())
		}

		r.log.Error("repository error", slog.Any("taskID", task.Id), slog.Any("userID", task.AuthorId), slog.Any("err", err))

		return nil, core_errors.ServerError()
	}

	r.log.Info("task successfully created in database", slog.Any("taskID", task.Id))

	domainTask := modelToDomain(saveTask)

	return &domainTask, nil

}
