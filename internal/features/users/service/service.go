package service

import (
	"Board_of_issuses/internal/core/domain"
	"context"
	"log/slog"
)

type UsersService struct {
	usersRepository UsersRepository
	authService     AuthService
	log             *slog.Logger
}

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_service.go -package=mocks

type UsersRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, email, login string) (*domain.User, error)
	ChangeUsersData(ctx context.Context, user *domain.User) (*domain.User, error)
}

type AuthService interface {
	CreateJWT(userID string) (string, error)
	GetUserIdFromJWT(JWT string) (string, error)
	ValidateJWT(JWT string) (*domain.Claims, error)
}

func NewUsersService(usersRepository UsersRepository, authService AuthService, log *slog.Logger) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		authService:     authService,
		log:             log,
	}
}
