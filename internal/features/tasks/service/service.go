package service

import (
	"Board_of_issuses/internal/core/domain"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type TasksService struct {
	tasksRepository TasksRepository
	log             *slog.Logger
}

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_service.go -package=mocks
type TasksRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	CompleteTask(ctx context.Context, userID, taskID uuid.UUID) error
	DeleteTask(ctx context.Context, taskID, authorID uuid.UUID) error

	UpdateTask(ctx context.Context, updateData *domain.UpdateTask) (*domain.Task, error)

	GetTasksFromOneDesk(ctx context.Context, filter *domain.TaskFilter) ([]*domain.Task, error)
	GetTaskByID(ctx context.Context, taskID, userID uuid.UUID) (*domain.Task, error)
}

func NewTasksService(tasksRepository TasksRepository, log *slog.Logger) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
		log:             log,
	}
}
