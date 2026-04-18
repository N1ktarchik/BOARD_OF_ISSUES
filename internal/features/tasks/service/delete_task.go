package service

import (
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *TasksService) DeleteTask(ctx context.Context, taskID, authorID string) error {
	s.log.Info("deleting task", slog.Any("task_id", taskID), slog.Any("author_id", authorID))

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		s.log.Warn("delete task failed: invalid taskID format", slog.Any("task_id", taskID), slog.Any("err", err))
		return core_errors.BadRequest()
	}

	if taskUUID == uuid.Nil {
		s.log.Warn("delete task failed: taskID is nil", slog.Any("task_id", taskID))
		return core_errors.BadRequest()
	}

	authorIDUUID, err := uuid.Parse(authorID)
	if err != nil {
		s.log.Warn("delete task failed: invalid authorID format", slog.Any("author_id", authorID), slog.Any("err", err))
		return core_errors.BadRequest()
	}

	if authorIDUUID == uuid.Nil {
		s.log.Warn("delete task failed: authorID is nil", slog.Any("author_id", authorID))
		return core_errors.BadRequest()
	}

	if err := s.tasksRepository.DeleteTask(ctx, taskUUID, authorIDUUID); err != nil {
		s.log.Error("delete task failed: repository error",
			slog.Any("task_id", taskID), slog.Any("author_id", authorID), slog.Any("err", err))

		return err
	}

	s.log.Info("task deleted successfully", slog.Any("task_id", taskID))
	return nil
}
