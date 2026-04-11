package handlers

import (
	dn "Board_of_issuses/internal/core/domains"
	"context"
)

type Service interface {
	Registration(ctx context.Context, user *dn.User) (string, error)
	Authorization(ctx context.Context, user *dn.User) (string, error)

	ChangeUserName(ctx context.Context, name string, userID int) error
	ChangeUserEmail(ctx context.Context, email string, userID int) error
	ChangeUserPassword(ctx context.Context, password string, userID int) error
	ConnectUserToDesk(ctx context.Context, userID, deskID int, password string) error

	CreateDesk(ctx context.Context, desk *dn.Desk) error
	ChangeDeskName(ctx context.Context, name string, deskId, userID int) error
	ChangeDeskPassword(ctx context.Context, password string, deskId, userID int) error
	ChangeDeskOwner(ctx context.Context, deskId, userID, newOwner int) error
	DeleteDesk(ctx context.Context, deskId, userID int) error
	GetAllDesks(ctx context.Context, userID int) ([]int, error)

	ComplyteTask(ctx context.Context, userID, taskID int) error
	UpdateTaskTime(ctx context.Context, userID, taskID, userTime int) error
	ChangeTaskDescription(ctx context.Context, userID, taskID int, description string) error
	DeleteTask(ctx context.Context, taskID, userID int) error
	CreateTask(ctx context.Context, task *dn.Task) error
	GetTasksWithParams(ctx context.Context, userId, deskID int, done bool) ([]dn.Task, error)
	GetAllTasks(ctx context.Context, userId, deskID int) ([]dn.Task, error)
}

type UserHandler struct {
	serv Service
}

func NewUserHandler(src Service) *UserHandler {
	return &UserHandler{
		serv: src,
	}
}
