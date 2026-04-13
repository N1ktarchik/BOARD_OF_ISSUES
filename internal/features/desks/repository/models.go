package repository

import (
	"Board_of_issuses/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type deskModel struct {
	Id         uuid.UUID
	Name       string
	Password   string
	OwnerId    uuid.UUID
	Created_at time.Time
}

func (m *deskModel) scan(row pgx.Row) error {
	return row.Scan(
		&m.Id,
		&m.Name,
		&m.Password,
		&m.OwnerId,
		&m.Created_at,
	)
}

func modelToDomain(model deskModel) domain.Desk {
	return domain.Desk{
		Id:         model.Id,
		Name:       model.Name,
		Password:   model.Password,
		OwnerId:    model.OwnerId,
		Created_at: model.Created_at,
	}
}
