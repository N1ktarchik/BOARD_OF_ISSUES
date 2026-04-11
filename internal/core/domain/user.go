package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" example:"636e856-e12b-56d9-f987-333222561234"`
	Login     string    `json:"login" example:"bonya123"`
	Password  string    `json:"password" example:"pass2000"`
	Email     string    `json:"email" example:"bonya123@example.com"`
	Name      string    `json:"name" example:"Bonya"`
	CreatedAt time.Time `json:"created_at" example:"2026-01-01T00:00:00Z"`
}
