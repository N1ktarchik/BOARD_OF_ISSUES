package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"N1ktarchik/Board_of_issues/internal/core/domain"
	"N1ktarchik/Board_of_issues/internal/features/tasks/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupService(t *testing.T) (*TasksService, *mocks.MockTasksRepository, func()) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTasksRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewTasksService(repo, logger)

	return svc, repo, func() { ctrl.Finish() }
}

func TestTasksService_CreateTask(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	task := &domain.Task{
		AuthorId: uuid.New(),
		DeskId:   uuid.New(),
		Name:     "Valid Task",
		Deadline: time.Now().Add(24 * time.Hour),
	}

	repo.EXPECT().
		CreateTask(ctx, gomock.Any()).
		Return(&domain.Task{Id: uuid.New()}, nil)

	result, err := svc.CreateTask(ctx, task)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestTasksService_CreateTask_NilAuthorID(t *testing.T) {
	svc, _, teardown := setupService(t)
	defer teardown()

	task := &domain.Task{
		DeskId:   uuid.New(),
		Name:     "Valid Task",
		Deadline: time.Now().Add(24 * time.Hour),
	}

	result, err := svc.CreateTask(context.Background(), task)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestTasksService_CreateTask_NameTooShort(t *testing.T) {
	svc, _, teardown := setupService(t)
	defer teardown()

	task := &domain.Task{
		AuthorId: uuid.New(),
		DeskId:   uuid.New(),
		Name:     "Ab", // Меньше 3 символов
		Deadline: time.Now().Add(24 * time.Hour),
	}

	result, err := svc.CreateTask(context.Background(), task)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestTasksService_CreateTask_RepoError(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	task := &domain.Task{
		AuthorId: uuid.New(),
		DeskId:   uuid.New(),
		Name:     "Valid Task",
		Deadline: time.Now().Add(24 * time.Hour),
	}

	expectedErr := errors.New("database connection lost")
	repo.EXPECT().
		CreateTask(ctx, gomock.Any()).
		Return(nil, expectedErr)

	result, err := svc.CreateTask(ctx, task)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Nil(t, result)
}

func TestTasksService_GetTaskByID(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	taskUUID := uuid.New()
	userUUID := uuid.New()

	repo.EXPECT().
		GetTaskByID(ctx, taskUUID, userUUID).
		Return(&domain.Task{Id: taskUUID}, nil)

	result, err := svc.GetTaskByID(ctx, taskUUID.String(), userUUID.String())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, taskUUID, result.Id)
}

func TestTasksService_GetTaskByID_InvalidUUID(t *testing.T) {
	svc, _, teardown := setupService(t)
	defer teardown()

	result, err := svc.GetTaskByID(context.Background(), "not-a-uuid", uuid.NewString())

	require.Error(t, err)
	require.Nil(t, result)
}

func TestTasksService_GetTaskByID_RepoError(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	taskUUID := uuid.New()
	userUUID := uuid.New()
	expectedErr := errors.New("task not found")

	repo.EXPECT().
		GetTaskByID(ctx, taskUUID, userUUID).
		Return(nil, expectedErr)

	result, err := svc.GetTaskByID(ctx, taskUUID.String(), userUUID.String())

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Nil(t, result)
}

func TestTasksService_DeleteTask(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	taskUUID := uuid.New()
	authorUUID := uuid.New()

	repo.EXPECT().
		DeleteTask(ctx, taskUUID, authorUUID).
		Return(nil)

	err := svc.DeleteTask(ctx, taskUUID.String(), authorUUID.String())

	require.NoError(t, err)
}

func TestTasksService_DeleteTask_InvalidAuthorID(t *testing.T) {
	svc, _, teardown := setupService(t)
	defer teardown()

	err := svc.DeleteTask(context.Background(), uuid.NewString(), "invalid-uuid")

	require.Error(t, err)
}

func TestTasksService_GetTasksFromOneDesk(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	deskUUID := uuid.New()
	userUUID := uuid.New()

	filter := &domain.TaskFilter{
		DeskId: deskUUID,
		UserId: userUUID,
	}

	expectedTasks := []*domain.Task{
		{Id: uuid.New(), Name: "Task 1", DeskId: deskUUID},
		{Id: uuid.New(), Name: "Task 2", DeskId: deskUUID},
	}

	repo.EXPECT().
		GetTasksFromOneDesk(ctx, filter).
		Return(expectedTasks, nil)

	result, err := svc.GetTasksFromOneDesk(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Equal(t, expectedTasks[0].Id, result[0].Id)
}

func TestTasksService_GetTasksFromOneDesk_NilDeskID(t *testing.T) {
	svc, _, teardown := setupService(t)
	defer teardown()

	filter := &domain.TaskFilter{
		DeskId: uuid.Nil,
		UserId: uuid.New(),
	}

	result, err := svc.GetTasksFromOneDesk(context.Background(), filter)

	require.Error(t, err)
	require.Nil(t, result)

}

func TestTasksService_GetTasksFromOneDesk_NilUserID(t *testing.T) {
	svc, _, teardown := setupService(t)
	defer teardown()

	filter := &domain.TaskFilter{
		DeskId: uuid.New(),
		UserId: uuid.Nil,
	}

	result, err := svc.GetTasksFromOneDesk(context.Background(), filter)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestTasksService_GetTasksFromOneDesk_RepoError(t *testing.T) {
	svc, repo, teardown := setupService(t)
	defer teardown()

	ctx := context.Background()
	filter := &domain.TaskFilter{
		DeskId: uuid.New(),
		UserId: uuid.New(),
	}

	expectedErr := errors.New("database is down")
	repo.EXPECT().
		GetTasksFromOneDesk(ctx, filter).
		Return(nil, expectedErr)

	result, err := svc.GetTasksFromOneDesk(ctx, filter)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedErr, err)
}
