package repository

import (
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *TasksRepository) CompleteTask(ctx context.Context, userID, taskID uuid.UUID) error {
	r.log.Info("completing task in repository", slog.Any("task_id", taskID), slog.Any("user_id", userID))

	query := `UPDATE tasks SET done=true WHERE id = $1 
			AND desk_id IN(
			SELECT desk_id FROM desk_members WHERE user_id = $2);`

	result, err := r.pool.Exec(ctx, query, taskID, userID)
	if err != nil {
		r.log.Error("repository error", slog.Any("userID", userID), slog.Any("err", err))
		return core_errors.ServerError()
	}

	if result.RowsAffected() == 0 {
		r.log.Warn("user have not access to desk", slog.Any("userID", userID), slog.Any("taskID", taskID))

		return core_errors.UserHaveNotAccessToDesk(userID.String(), "")
	}

	r.log.Info("task completed successfully in repository",
		slog.Any("task_id", taskID),
		slog.Any("user_id", userID))

	return nil
}
