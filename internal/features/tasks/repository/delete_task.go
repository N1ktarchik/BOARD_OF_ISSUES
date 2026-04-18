package repository

import (
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, taskID, authorID uuid.UUID) error {
	r.log.Info("deleting task in repository", slog.Any("task_id", taskID), slog.Any("author_id", authorID))

	query := `DELETE FROM tasks WHERE id=$1 AND author_id=$2`

	result, err := r.pool.Exec(ctx, query, taskID, authorID)
	if err != nil {
		r.log.Error("repository error", slog.Any("userID", authorID), slog.Any("err", err))
		return core_errors.ServerError()
	}

	if result.RowsAffected() == 0 {
		r.log.Warn("user is not owner of a task or task not found",
			slog.Any("authorID", authorID),
			slog.Any("taskID", taskID))

		return core_errors.UserNotOwnerOfTask(authorID.String(), taskID.String())
	}

	r.log.Info("task deleted successfully in repository", slog.Any("task_id", taskID))
	return nil
}
