package http

import (
	dn "Board_of_issuses/internal/core/domain"

	"github.com/google/uuid"
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

type ConnectUserToDeskRequestDTO struct {
	ID       uuid.UUID `json:"desk_id" example:"636e856-e12b-56d9-f987-333222561234"`
	Password string    `json:"password" example:"deskpass2000"`
}
