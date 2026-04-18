package service

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (s *TasksService) UpdateTask(ctx context.Context, updateData *domain.UpdateTask) (*domain.Task, error) {
	s.log.Info("updating task", slog.Any("task_id", updateData.Id))

	if updateData.Id == uuid.Nil {
		s.log.Warn("update task failed: taskID is nil", slog.Any("task_id", updateData.Id))
		return nil, core_errors.BadRequest()
	}

	if updateData.AuthorId == uuid.Nil {
		s.log.Warn("update task failed: authorID is nil", slog.Any("author_id", updateData.AuthorId))
		return nil, core_errors.BadRequest()
	}

	if updateData.Name != "" && len(updateData.Name) < 3 {
		s.log.Warn("update task failed: name length is less than 3", slog.String("task name", updateData.Name))
		return nil, core_errors.BadRequest()
	}

	if !updateData.Deadline.IsZero() && updateData.Deadline.UTC().Before(time.Now().UTC()) {
		s.log.Warn("update task failed: deadline is in the past", slog.Any("deadline", updateData.Deadline))
		return nil, core_errors.BadRequest()
	}

	updateTask, err := s.tasksRepository.UpdateTask(ctx, updateData)
	if err != nil {
		s.log.Error("update task failed: repository error", slog.Any("task_id", updateData.Id), slog.Any("err", err))
		return nil, err
	}

	s.log.Info("task updated successfully", slog.Any("task_id", updateData.Id))
	return updateTask, nil

}
