package service

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"
)

func (s *UsersService) LoginUser(ctx context.Context, user *domain.User) (string, error) {

	if user.Password == "" || (user.Email == "" && user.Login == "") {
		s.log.Warn("login failed: empty credentials")

		return "", core_errors.BadRequest()
	}

	repoUser, err := s.usersRepository.GetUser(ctx, user.Login, user.Email)
	if err != nil {
		s.log.Error("login failed: user not found", slog.Any("err", err))
		return "", err
	}

	if !domain.Compare(user.Password, repoUser.Password) {
		s.log.Warn("login failed: wrong password")
		return "", core_errors.InvalidPassword()
	}

	JWTtoken, err := s.authService.CreateJWT(repoUser.ID.String())
	if err != nil {
		s.log.Error("login failed: jwt generation error", slog.Any("err", err))

		return "", err
	}

	s.log.Info("user logged in successfully", slog.String("user_id", repoUser.ID.String()))
	return JWTtoken, nil

}
