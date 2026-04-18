package http

import (
	"Board_of_issuses/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type DeskRequestDTO struct {
	Id         uuid.UUID `json:"id" example:"636e856-e12b-56d9-f987-333222561234"`
	Name       string    `json:"name" example:"My Desk"`
	Password   string    `json:"password" example:"mysecretpassword"`
	OwnerId    uuid.UUID `json:"owner_id" example:"832t758-a12g-47y9-i999-123456789098"`
	Created_at time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}


func (d *DeskRequestDTO) ToServiceDesk() *domain.Desk {
	return &domain.Desk{
		Id:         d.Id,
		Name:       d.Name,
		Password:   d.Password,
		OwnerId:    d.OwnerId,
		Created_at: d.Created_at,
	}
}

