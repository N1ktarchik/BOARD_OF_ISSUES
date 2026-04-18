package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type TaskRequestDTO struct {
	Id          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AuthorId    uuid.UUID `json:"author_id" example:"636e856-e12b-56d9-f987-333222561234"`
	DeskId      uuid.UUID `json:"desk_id" example:"832t758-a12g-47y9-i999-123456789098"`
	Name        string    `json:"name" example:"Task name"`
	Description string    `json:"description" example:"Task description"`
	Done        bool      `json:"status" example:"false"`
	Deadline    time.Time `json:"deadline" example:"2023-10-10T10:00:00Z"`
	Created_at  time.Time `json:"created_at" example:"2023-10-10T10:00:00Z"`
}

type UpdateTaskRequestDTO struct {
	Id          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AuthorId    uuid.UUID `json:"author_id" example:"636e856-e12b-56d9-f987-333222561234"`
	Name        string    `json:"name" example:"Task name"`
	Description string    `json:"description" example:"Task description"`
	Deadline    time.Time `json:"deadline" example:"2023-10-10T10:00:00Z"`
}

func (u *UpdateTaskRequestDTO) ToServiceUpdateTask() *domain.UpdateTask {
	return &domain.UpdateTask{
		Id:          u.Id,
		AuthorId:    u.AuthorId,
		Name:        u.Name,
		Description: u.Description,
		Deadline:    u.Deadline,
	}
}

func (t *TaskRequestDTO) ToServiceTask() *domain.Task {
	return &domain.Task{
		Id:          t.Id,
		AuthorId:    t.AuthorId,
		DeskId:      t.DeskId,
		Name:        t.Name,
		Description: t.Description,
		Done:        t.Done,
		Deadline:    t.Deadline,
		Created_at:  t.Created_at,
	}
}
