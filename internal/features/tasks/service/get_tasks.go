package service

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *TasksService) GetTasksFromOneDesk(ctx context.Context, filter *domain.TaskFilter) ([]*domain.Task, error) {
	s.log.Info("getting tasks", slog.Any("filter", filter))

	if filter.DeskId == uuid.Nil {
		s.log.Warn("get tasks failed: invalid deskID", slog.Any("desk_id", filter.DeskId))
		return nil, core_errors.BadRequest()
	}

	if filter.UserId == uuid.Nil {
		s.log.Warn("get tasks failed: invalid userID", slog.Any("user_id", filter.UserId))
		return nil, core_errors.BadRequest()
	}

	tasks, err := s.tasksRepository.GetTasksFromOneDesk(ctx, filter)
	if err != nil {
		s.log.Error("get tasks failed: repository error", slog.Any("filter", filter), slog.Any("err", err))
		return nil, err
	}

	s.log.Info("tasks retrieved successfully", slog.Any("filter", filter))
	return tasks, nil

}

func (s *TasksService) GetTaskByID(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	s.log.Info("getting task by id", slog.Any("task_id", taskID), slog.Any("user_id", userID))

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		s.log.Warn("get task by id failed: error parsing task id", slog.Any("task_id", taskID), slog.Any("err", err))
		return nil, core_errors.BadRequest()
	}

	if taskUUID == uuid.Nil {
		s.log.Warn("get task by id failed: empty task id")
		return nil, core_errors.BadRequest()
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.log.Warn("get task by id failed: error parsing user id", slog.Any("user_id", userID), slog.Any("err", err))
		return nil, core_errors.BadRequest()
	}

	if userUUID == uuid.Nil {
		s.log.Warn("get task by id failed: empty user id")
		return nil, core_errors.BadRequest()
	}

	task, err := s.tasksRepository.GetTaskByID(ctx, taskUUID, userUUID)
	if err != nil {
		s.log.Error("get task by id failed: repository error",
			slog.Any("task_id", taskID), slog.Any("user_id", userID), slog.Any("err", err))

		return nil, err
	}

	s.log.Info("task retrieved successfully", slog.Any("task_id", taskID), slog.Any("user_id", userID))
	return task, nil

}
