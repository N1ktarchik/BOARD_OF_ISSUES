package service

import (
	"Board_of_issuses/internal/core/domain"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type DesksService struct {
	deskRepository DeskRepository
	log            *slog.Logger
}

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_service.go -package=mocks

type DeskRepository interface {
	CreateDesk(ctx context.Context, desk *domain.Desk) (*domain.Desk, error)
	ChangeDesksData(ctx context.Context, deskUpdate *domain.Desk, requesterID uuid.UUID) (*domain.Desk, error)
	DeleteDesk(ctx context.Context, userUUID, deskUUID uuid.UUID) error

	ConnectUserToDesk(ctx context.Context, userID, deskID uuid.UUID) error

	GetAllUsersDesks(ctx context.Context, userUUID uuid.UUID) ([]domain.Desk, error)
}

func NewDesksService(deskRepository DeskRepository, log *slog.Logger) *DesksService {
	return &DesksService{
		deskRepository: deskRepository,
		log:            log,
	}
}
