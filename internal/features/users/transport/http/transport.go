package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"context"
	"log/slog"
)

type UsersHandler struct {
	usersService UsersService
	log          *slog.Logger
}

type UsersService interface {
	RegisterUser(ctx context.Context, user *domain.User) (string, error)
	LoginUser(ctx context.Context, user *domain.User) (string, error)

	ChangeUsersData(ctx context.Context, user *domain.User) (*domain.User, error)
}

func NewUsersHandler(usersService UsersService, log *slog.Logger) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
		log:          log,
	}
}
