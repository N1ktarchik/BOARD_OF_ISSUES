package service

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"
	"strings"
)

func (s *UsersService) ChangeUsersData(ctx context.Context, user *domain.User) (*domain.User, error) {

	if user.Name == "" && user.Email == "" && user.Password == "" {
		s.log.Warn("change user data successfully : no data to change", slog.String("user_id", user.ID.String()))
		return user, nil
	}

	if user.Name != "" && len(user.Name) < 3 {
		s.log.Warn("change users data failed : invalid user name", slog.String("user_id", user.ID.String()))
		return nil, core_errors.BadRequest()
	}

	if user.Email != "" {

		if len(user.Email) < 7 || !strings.Contains(user.Email, "@") {
			s.log.Warn("change users data failed : invalid user email", slog.String("user_id", user.ID.String()))
			return nil, core_errors.BadRequest()
		}
	}

	hashPassword, err := domain.Hash(user.Password)
	if err != nil {
		s.log.Warn("change users data failed : password hashing error", slog.Any("err", err),
			slog.String("user_id", user.ID.String()))

		return nil, err
	}

	user.Password = hashPassword

	saveUser, err := s.usersRepository.ChangeUsersData(ctx, user)
	if err != nil {
		s.log.Error("change users data failed : repository error", slog.Any("err", err))

		return nil, err
	}

	s.log.Info("users data changed successfully", slog.String("user_id", user.ID.String()))
	return saveUser, nil

}
