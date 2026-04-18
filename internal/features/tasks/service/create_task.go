package service

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (s *TasksService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	s.log.Info("creating task", slog.Any("authorID", task.AuthorId), slog.Any("deskID", task.DeskId))

	if task.AuthorId == uuid.Nil {
		s.log.Warn("create task failed: authorID is nil")
		return nil, core_errors.BadRequest()
	}

	if task.DeskId == uuid.Nil {
		s.log.Warn("create task failed: deskID is nil")
		return nil, core_errors.BadRequest()
	}

	if len(task.Name) < 3 {
		s.log.Warn("create task failed: name length is less than 3", slog.String("task name", task.Name))
		return nil, core_errors.BadRequest()
	}

	if task.Deadline.UTC().Before(time.Now().UTC()) {
		s.log.Warn("create task failed: deadline is in the past", slog.Any("deadline", task.Deadline))
		return nil, core_errors.BadRequest()
	}

	if task.Deadline.IsZero() {
		s.log.Warn("task deadline is zero, setting default deadline to 30 days from now")
		task.Deadline = task.Deadline.Add(30 * 24 * time.Hour)
	}

	task.Id = uuid.New()

	saveTask, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		s.log.Error("create task failed: repository error", slog.Any("authorID", task.AuthorId), slog.Any("err", err))
		return nil, err
	}

	s.log.Info("task created successfully", slog.Any("taskID", saveTask.Id))
	return saveTask, nil
}
