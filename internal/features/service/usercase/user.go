package usercase

import (
	er "Board_of_issuses/internal/core"
	"Board_of_issuses/internal/core/auth"
	dn "Board_of_issuses/internal/core/domains"
	"context"
)

func (s *Service) Registration(ctx context.Context, user *dn.User) (string, error) {

	register, err := s.repo.CheckUserByEmailOrLogin(ctx, user.Login, user.Email)
	if err != nil {
		return "", err
	}
	if register {
		return "", er.HaveRegister(user.Login, user.Email)
	}

	hashPassword, err := auth.Hash(user.Password)
	if err != nil {
		return "", err
	}

	user.Password = hashPassword
	userId, err := s.repo.CreateUser(ctx, user.ToRepoUser())
	if err != nil {
		return "", err
	}

	user.Id = userId

	token, err := s.auth.JWTManager.Create(user.Id, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Authorization(ctx context.Context, user *dn.User) (string, error) {

	repoUser, err := s.repo.GetUserByLoginOrEmail(ctx, user.Login, user.Email)
	if err != nil {
		return "", err
	}

	if !auth.Compare(user.Password, repoUser.Password) {
		return "", er.InvalidPassword()
	}

	token, err := s.auth.JWTManager.Create(repoUser.Id, repoUser.Email)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *Service) ChangeUserName(ctx context.Context, name string, userID int) error {
	if err := s.repo.UpdateUserName(ctx, name, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ChangeUserEmail(ctx context.Context, email string, userID int) error {
	if err := s.repo.UpdateUserEmail(ctx, email, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ChangeUserPassword(ctx context.Context, password string, userID int) error {

	hashPassword, err := auth.Hash(password)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateUserPassword(ctx, hashPassword, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ConnectUserToDesk(ctx context.Context, userID, deskID int, password string) error {
	validDeskPassword, err := s.repo.CheckDeskPassword(ctx, deskID)
	if err != nil {
		return err
	}

	if !auth.Compare(password, validDeskPassword) {
		return er.InvalidPassword()
	}

	if err := s.repo.ConnectUserToDesk(ctx, userID, deskID); err != nil {
		return err
	}

	return nil
}
