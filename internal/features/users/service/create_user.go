package service

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *UsersService) RegisterUser(ctx context.Context, user *domain.User) (string, error) {

	if user.Login == "" || user.Password == "" || user.Email == "" || user.Name == "" {
		s.log.Error("register failed: invalid user data")
		return "", core_errors.BadRequest()
	}

	hashPassword, err := domain.Hash(user.Password)
	if err != nil {
		s.log.Error("register failed: password hashing error", slog.Any("err", err))
		return "", err
	}

	user.Password = hashPassword
	user.ID = uuid.New()

	if err := s.usersRepository.CreateUser(ctx, user); err != nil {
		s.log.Error("register failed: repository error", slog.Any("err", err))
		return "", err
	}

	JWTtoken, err := s.authService.CreateJWT(user.ID.String())
	if err != nil {
		s.log.Error("register failed: jwt generation error", slog.Any("err", err))
		return "", err
	}

	s.log.Info("user registered successfully", slog.String("user_id", user.ID.String()))

	return JWTtoken, nil
}
