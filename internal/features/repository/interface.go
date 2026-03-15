package repository

import (
	"context"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (int, error)
	UpdateUserEmail(ctx context.Context, email string, userId int) error
	UpdateUserName(ctx context.Context, name string, userId int) error
	UpdateUserPassword(ctx context.Context, password string, userId int) error
	DeleteUser(ctx context.Context, userId string) error
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByLoginOrEmail(ctx context.Context, login, email string) (*User, error)
	CheckUserByEmailOrLogin(ctx context.Context, login, email string) (bool, error)

	ConnectUserToDesk(ctx context.Context, userID, deskID int) error
	GetUserDesks(ctx context.Context, userId int) ([]int, error)
	DeleteUserDesk(ctx context.Context, userId, deskId int) error
	CheckUserDesk(ctx context.Context, userId, deskId int) (bool, error)

	CreateDesk(ctx context.Context, desk *Desk) error
	UpdateDeskName(ctx context.Context, deskId int, name string) error
	UpdateDesksPassword(ctx context.Context, deskId int, password string) error
	UpdateDeskOwner(ctx context.Context, ownerid, deskid int) error
	DeleteDesk(ctx context.Context, deskId int) error
	CheckDeskOwner(ctx context.Context, deskId int) (int, error)
	CheckDeskPassword(ctx context.Context, deskId int) (string, error)

	CreateTask(ctx context.Context, task *Task) error
	UpdateTaskDecription(ctx context.Context, id int, description string) error
	UpdateTaskTime(ctx context.Context, id int, time time.Time) error
	UpdateTaskDone(ctx context.Context, id int) error
	DeleteTask(ctx context.Context, id int) error
	GetAllTasksFromOneDesk(ctx context.Context, deskId int) ([]Task, error)
	GetTasksWithParams(ctx context.Context, deskId int, done bool) ([]Task, error)
	GetTaskOwner(ctx context.Context, taskID int) (int, error)
	GetDeskIDByTask(ctx context.Context, taskID int) (int, error)
}
