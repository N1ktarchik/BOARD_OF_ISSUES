package http

import (
	"Board_of_issuses/internal/core/domain"
	"context"
	"log/slog"
)

type TasksHandler struct {
	tasksService TasksService
	log          *slog.Logger
}

type TasksService interface {
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, taskID, authorID string) error
	CompleteTask(ctx context.Context, userID, taskID string) error

	UpdateTask(ctx context.Context, updateData *domain.UpdateTask) (*domain.Task, error)

	GetTasksFromOneDesk(ctx context.Context, filter *domain.TaskFilter) ([]*domain.Task, error)
	GetTaskByID(ctx context.Context, taskID, userID string) (*domain.Task, error)
}

func NewTasksHandler(tasksService TasksService, log *slog.Logger) *TasksHandler {
	return &TasksHandler{
		tasksService: tasksService,
		log:          log,
	}
}
