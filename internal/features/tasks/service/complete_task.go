package service

import (
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *TasksService) CompleteTask(ctx context.Context, userID, taskID string) error {
	s.log.Info("completing task", slog.Any("task_id", taskID), slog.Any("user_id", userID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.log.Warn("complete task failed: invalid userID format", slog.Any("user_id", userID), slog.Any("err", err))
		return core_errors.BadRequest()
	}

	if userUUID == uuid.Nil {
		s.log.Warn("complete task failed: userID is nil", slog.Any("user_id", userID))
		return core_errors.BadRequest()
	}

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		s.log.Warn("complete task failed: invalid taskID format", slog.Any("task_id", taskID), slog.Any("err", err))
		return core_errors.BadRequest()
	}

	if taskUUID == uuid.Nil {
		s.log.Warn("complete task failed: taskID is nil", slog.Any("task_id", taskID))
		return core_errors.BadRequest()
	}

	if err := s.tasksRepository.CompleteTask(ctx, userUUID, taskUUID); err != nil {
		s.log.Error("complete task failed: repository error",
			slog.Any("task_id", taskID), slog.Any("user_id", userID), slog.Any("err", err))

		return err
	}

	return nil

}
