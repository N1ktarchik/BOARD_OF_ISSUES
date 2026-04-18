package service

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"N1ktarchik/Board_of_issues/internal/features/desks/service/mocks"
)

func initTest(t *testing.T) (*DesksService, *mocks.MockDeskRepository, context.Context) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockDeskRepository(ctrl)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewDesksService(repo, log)
	return svc, repo, context.Background()
}

func TestGetAllUsersDesks_Success(t *testing.T) {
	svc, repo, ctx := initTest(t)
	userID := uuid.New()

	expectedDesks := []domain.Desk{
		{Id: uuid.New(), Name: "Desk 1"},
		{Id: uuid.New(), Name: "Desk 2"},
	}

	repo.EXPECT().GetAllUsersDesks(ctx, userID).Return(expectedDesks, nil)

	desks, err := svc.GetAllUsersDesks(ctx, userID.String())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(desks) != 2 {
		t.Errorf("expected 2 desks, got %d", len(desks))
	}
}

func TestGetAllUsersDesks_InvalidUUID(t *testing.T) {
	svc, _, ctx := initTest(t)

	_, err := svc.GetAllUsersDesks(ctx, "invalid-uuid")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAllUsersDesks_RepoError(t *testing.T) {
	svc, repo, ctx := initTest(t)
	userID := uuid.New()

	repo.EXPECT().GetAllUsersDesks(ctx, userID).Return(nil, errors.New("repo error"))

	_, err := svc.GetAllUsersDesks(ctx, userID.String())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeleteDesk_Success(t *testing.T) {
	svc, repo, ctx := initTest(t)
	userID := uuid.New()
	deskID := uuid.New()

	repo.EXPECT().DeleteDesk(ctx, userID, deskID).Return(nil)

	err := svc.DeleteDesk(ctx, deskID.String(), userID.String())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDeleteDesk_InvalidIDs(t *testing.T) {
	svc, _, ctx := initTest(t)
	validID := uuid.New().String()

	err1 := svc.DeleteDesk(ctx, validID, "invalid-user")
	if err1 == nil {
		t.Error("expected error for invalid user id")
	}

	err2 := svc.DeleteDesk(ctx, "invalid-desk", validID)
	if err2 == nil {
		t.Error("expected error for invalid desk id")
	}
}

func TestDeleteDesk_RepoError(t *testing.T) {
	svc, repo, ctx := initTest(t)
	userID := uuid.New()
	deskID := uuid.New()

	repo.EXPECT().DeleteDesk(ctx, userID, deskID).Return(errors.New("db error"))

	err := svc.DeleteDesk(ctx, deskID.String(), userID.String())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateDesk_Success(t *testing.T) {
	svc, repo, ctx := initTest(t)

	desk := &domain.Desk{
		Name:     "Test Desk",
		Password: "password123",
		OwnerId:  uuid.New(),
	}

	repo.EXPECT().CreateDesk(ctx, gomock.Any()).Return(&domain.Desk{Id: uuid.New(), Name: "Test Desk"}, nil)

	created, err := svc.CreateDesk(ctx, desk)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created == nil {
		t.Fatal("expected desk, got nil")
	}
}

func TestCreateDesk_ShortName(t *testing.T) {
	svc, _, ctx := initTest(t)

	desk := &domain.Desk{
		Name:     "ab",
		Password: "password123",
		OwnerId:  uuid.New(),
	}

	_, err := svc.CreateDesk(ctx, desk)

	if err == nil {
		t.Error("expected error for short name, got nil")
	}
}

func TestCreateDesk_EmptyOwnerID(t *testing.T) {
	svc, _, ctx := initTest(t)

	desk := &domain.Desk{
		Name:     "Test Desk",
		Password: "password123",
		OwnerId:  uuid.Nil,
	}

	_, err := svc.CreateDesk(ctx, desk)

	if err == nil {
		t.Error("expected error for empty owner, got nil")
	}
}

func TestConnectUserToDesk_Success(t *testing.T) {
	svc, repo, ctx := initTest(t)
	userID := uuid.New()
	deskID := uuid.New()

	repo.EXPECT().ConnectUserToDesk(ctx, userID, deskID).Return(nil)

	err := svc.ConnectUserToDesk(ctx, userID, deskID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestConnectUserToDesk_EmptyIDs(t *testing.T) {
	svc, _, ctx := initTest(t)
	validID := uuid.New()

	err1 := svc.ConnectUserToDesk(ctx, uuid.Nil, validID)
	if err1 == nil {
		t.Error("expected error for empty user id")
	}

	err2 := svc.ConnectUserToDesk(ctx, validID, uuid.Nil)
	if err2 == nil {
		t.Error("expected error for empty desk id")
	}
}

func TestChangeDesksData_Success(t *testing.T) {
	svc, repo, ctx := initTest(t)
	requesterID := uuid.New()

	deskUpdate := &domain.Desk{
		Id:       uuid.New(),
		Name:     "New Desk Name",
		Password: "newpassword123",
		OwnerId:  requesterID,
	}

	repo.EXPECT().ChangeDesksData(ctx, gomock.Any(), requesterID).Return(&domain.Desk{Id: deskUpdate.Id, Name: "New Desk Name"}, nil)

	res, err := svc.ChangeDesksData(ctx, deskUpdate, requesterID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil {
		t.Fatal("expected updated desk, got nil")
	}
}

func TestChangeDesksData_ShortName(t *testing.T) {
	svc, _, ctx := initTest(t)
	requesterID := uuid.New()

	deskUpdate := &domain.Desk{
		Id:      uuid.New(),
		Name:    "no",
		OwnerId: requesterID,
	}

	_, err := svc.ChangeDesksData(ctx, deskUpdate, requesterID)

	if err == nil {
		t.Error("expected error for short name")
	}
}

func TestChangeDesksData_EmptyRequesterID(t *testing.T) {
	svc, _, ctx := initTest(t)

	deskUpdate := &domain.Desk{
		Id:      uuid.New(),
		Name:    "Valid Name",
		OwnerId: uuid.New(),
	}

	_, err := svc.ChangeDesksData(ctx, deskUpdate, uuid.Nil)

	if err == nil {
		t.Error("expected error for empty requester id")
	}
}

func TestChangeDesksData_RepoError(t *testing.T) {
	svc, repo, ctx := initTest(t)
	requesterID := uuid.New()

	deskUpdate := &domain.Desk{
		Id:      uuid.New(),
		Name:    "Valid Name",
		OwnerId: requesterID,
	}

	repo.EXPECT().ChangeDesksData(ctx, gomock.Any(), requesterID).Return(nil, errors.New("update failed"))

	_, err := svc.ChangeDesksData(ctx, deskUpdate, requesterID)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
