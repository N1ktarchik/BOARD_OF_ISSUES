package repository

import (
	"Board_of_issuses/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userModel struct {
	ID        uuid.UUID
	Login     string
	Password  string
	Email     string
	Name      string
	CreatedAt time.Time
}

func (m *userModel) scan(row pgx.Row) error {
	return row.Scan(
		&m.ID,
		&m.Login,
		&m.Password,
		&m.Email,
		&m.Name,
		&m.CreatedAt,
	)
}

func modelToDomain(model userModel) domain.User {
	return domain.User{
		ID:        model.ID,
		Login:     model.Login,
		Password:  model.Password,
		Email:     model.Email,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
	}
}
