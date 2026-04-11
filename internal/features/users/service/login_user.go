package service

import (
	"Board_of_issuses/internal/core/domain"
	"Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"
)

func (s *UsersService) LoginUser(ctx context.Context, user *domain.User) (string, error) {

	if user.Password == "" || (user.Email == "" && user.Login == "") {
		s.log.Error("login failed: empty credentials")

		return "", errors.BadRequest()
	}

	repoUser, err := s.usersRepository.GetUser(ctx, user.Login, user.Email)
	if err != nil {
		s.log.Error("login failed: user not found", slog.Any("err", err))
		return "", err
	}

	if !domain.Compare(user.Password, repoUser.Password) {
		s.log.Error("login failed: wrong password")
		return "", errors.InvalidPassword()
	}

	JWTtoken, err := s.authService.CreateJWT(repoUser.ID.String())
	if err != nil {
		s.log.Error("login failed: jwt generation error", slog.Any("err", err))

		return "", err
	}

	s.log.Info("user logged in successfully", slog.String("user_id", repoUser.ID.String()))
	return JWTtoken, nil

}
