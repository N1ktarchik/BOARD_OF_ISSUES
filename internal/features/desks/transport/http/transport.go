package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type DesksHandler struct {
	desksService DesksService
	log          *slog.Logger
}

type DesksService interface {
	CreateDesk(ctx context.Context, desk *domain.Desk) (*domain.Desk, error)
	ChangeDesksData(ctx context.Context, deskUpdate *domain.Desk, requesterID uuid.UUID) (*domain.Desk, error)
	DeleteDesk(ctx context.Context, deskID, userID string) error

	GetAllUsersDesks(ctx context.Context, userID string) ([]domain.Desk, error)

	ConnectUserToDesk(ctx context.Context, userID, deskID uuid.UUID) error
}

func NewDesksHandler(desksService DesksService, log *slog.Logger) *DesksHandler {
	return &DesksHandler{
		desksService: desksService,
		log:          log,
	}
}
