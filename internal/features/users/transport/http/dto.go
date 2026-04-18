package http

import (
	dn "N1ktarchik/Board_of_issues/internal/core/domain"
)

type UsersRequestDTO struct {
	Login    string `json:"login" example:"bonya123"`
	Password string `json:"password" example:"pass2000"`
	Email    string `json:"email" example:"bonya123@example.com"`
	Name     string `json:"name" example:"Bonya"`
}

func (u *UsersRequestDTO) ToServiceUser() *dn.User {
	return &dn.User{
		Login:    u.Login,
		Password: u.Password,
		Email:    u.Email,
		Name:     u.Name,
	}
}
