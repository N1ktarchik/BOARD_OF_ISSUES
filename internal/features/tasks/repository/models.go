package repository

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type taskModel struct {
	Id          uuid.UUID
	AuthorId    uuid.UUID
	DeskId      uuid.UUID
	Name        string
	Description string
	Done        bool
	Deadline    time.Time
	Created_at  time.Time
}

func (m *taskModel) scan(row pgx.Row) error {
	return row.Scan(
		&m.Id,
		&m.AuthorId,
		&m.DeskId,
		&m.Name,
		&m.Description,
		&m.Done,
		&m.Deadline,
		&m.Created_at,
	)
}

func modelToDomain(model taskModel) domain.Task {
	return domain.Task{
		Id:          model.Id,
		AuthorId:    model.AuthorId,
		DeskId:      model.DeskId,
		Name:        model.Name,
		Description: model.Description,
		Done:        model.Done,
		Deadline:    model.Deadline,
		Created_at:  model.Created_at,
	}
}
