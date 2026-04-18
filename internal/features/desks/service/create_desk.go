package service

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"log/slog"

	"context"

	"github.com/google/uuid"
)

func (s *DesksService) CreateDesk(ctx context.Context, desk *domain.Desk) (*domain.Desk, error) {
	s.log.Info("creating desk", slog.Any("deskName", desk.Name), slog.Any("ownerID", desk.OwnerId))

	if len(desk.Name) < 3 {
		s.log.Warn("create desk failed: short desk name", slog.Any("deskName", desk.Name))
		return nil, core_errors.BadRequest()
	}

	if desk.OwnerId == uuid.Nil {
		s.log.Warn("create desk failed: empty owner id")
		return nil, core_errors.BadRequest()
	}

	hashPassword, err := domain.Hash(desk.Password)
	if err != nil {
		s.log.Warn("create desk failed: password hashing error", slog.Any("err", err))
		return nil, err
	}

	desk.Password = hashPassword
	desk.Id = uuid.New()

	saveDesk, err := s.deskRepository.CreateDesk(ctx, desk)
	if err != nil {
		s.log.Error("create desk failed: repository error",
			slog.Any("deskID", desk.Id), slog.Any("ownerID", desk.OwnerId), slog.Any("err", err))

		return nil, err
	}

	s.log.Info("desk created successfully", slog.Any("deskID", saveDesk.Id))
	return saveDesk, nil
}
