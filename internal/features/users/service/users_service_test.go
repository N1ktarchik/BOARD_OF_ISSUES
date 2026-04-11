package service

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"Board_of_issuses/internal/core/domain"
	"Board_of_issuses/internal/core/errors"
	"Board_of_issuses/internal/features/users/service/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func initTest(t *testing.T) (*UsersService, *mocks.MockUsersRepository, *mocks.MockAuthService, context.Context) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockUsersRepository(ctrl)
	auth := mocks.NewMockAuthService(ctrl)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := NewUsersService(repo, auth, logger)

	ctx := context.Background()

	return svc, repo, auth, ctx
}

func TestRegisterUser_Success(t *testing.T) {
	svc, repo, auth, ctx := initTest(t)
	user := &domain.User{Login: "user1", Password: "password123", Email: "test@mail.ru", Name: "Ivan"}

	repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil)
	auth.EXPECT().CreateJWT(gomock.Any()).Return("token-123", nil)

	token, err := svc.RegisterUser(ctx, user)
	if err != nil {
		t.Fatalf("Expected success, got err: %v", err)
	}
	if token != "token-123" {
		t.Errorf("Expected token token-123, got %s", token)
	}
}

func TestRegisterUser_AlreadyExists(t *testing.T) {
	svc, repo, _, ctx := initTest(t)
	user := &domain.User{Login: "busy_man", Password: "password123", Email: "a@a.ru", Name: "A"}

	repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(errors.UserAlreadyRegistered("busy_man", ""))

	_, err := svc.RegisterUser(ctx, user)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRegisterUser_EmptyFields(t *testing.T) {
	svc, _, _, ctx := initTest(t)
	user := &domain.User{Login: "", Password: "123"}

	_, err := svc.RegisterUser(ctx, user)
	if err == nil {
		t.Error("Expected error for empty fields, got nil")
	}
}

// Успешный вход
func TestLoginUser_Success(t *testing.T) {
	svc, repo, auth, ctx := initTest(t)

	rawPass := "password"
	hash, _ := domain.Hash(rawPass)
	id := uuid.New()

	repo.EXPECT().GetUser(ctx, "login", "").Return(&domain.User{ID: id, Password: hash}, nil)
	auth.EXPECT().CreateJWT(id.String()).Return("jwt-token", nil)

	token, err := svc.LoginUser(ctx, &domain.User{Login: "login", Password: rawPass})
	if err != nil || token != "jwt-token" {
		t.Errorf("Login failed: %v", err)
	}
}
func TestLoginUser_NotFound(t *testing.T) {
	svc, repo, _, ctx := initTest(t)

	repo.EXPECT().GetUser(ctx, "ghost", "").Return(nil, errors.ServerError()) // Или твоя ошибка Not Found

	_, err := svc.LoginUser(ctx, &domain.User{Login: "ghost", Password: "123"})
	if err == nil {
		t.Error("Expected error when user not found, got nil")
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	svc, repo, _, ctx := initTest(t)

	hash, _ := domain.Hash("correct")
	repo.EXPECT().GetUser(ctx, "user", "").Return(&domain.User{Password: hash}, nil)

	_, err := svc.LoginUser(ctx, &domain.User{Login: "user", Password: "WRONG"})
	if err == nil {
		t.Error("Expected 'InvalidPassword' error, got nil")
	}
}

func TestChangeData_ShortName(t *testing.T) {
	svc, _, _, ctx := initTest(t)
	user := &domain.User{ID: uuid.New(), Name: "Jo"}

	_, err := svc.ChangeUsersData(ctx, user)
	if err == nil {
		t.Error("Expected error for name < 3 chars, got nil")
	}
}

// Ошибка: кривой email
func TestChangeData_InvalidEmail(t *testing.T) {
	svc, _, _, ctx := initTest(t)
	user := &domain.User{ID: uuid.New(), Email: "not_an_email"}

	_, err := svc.ChangeUsersData(ctx, user)
	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
}

func TestChangeData_PasswordOnly(t *testing.T) {
	svc, repo, _, ctx := initTest(t)
	id := uuid.New()
	user := &domain.User{ID: id, Password: "new_password"}

	repo.EXPECT().ChangeUsersData(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *domain.User) (*domain.User, error) {
		if domain.Compare("new_password", u.Password) {
			return u, nil
		}
		return nil, io.ErrUnexpectedEOF
	})

	res, err := svc.ChangeUsersData(ctx, user)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
	if res.ID != id {
		t.Error("Returned user ID doesn't match")
	}
}
