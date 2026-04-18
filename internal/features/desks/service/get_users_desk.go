package service

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *DesksService) GetAllUsersDesks(ctx context.Context, userID string) ([]domain.Desk, error) {
	s.log.Info("getting all users desks", slog.Any("userID", userID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.log.Warn("get all users desks failed: error parsing user id", slog.Any("err", err))
		return nil, core_errors.BadRequest()
	}
	if userUUID == uuid.Nil {
		s.log.Warn("get all users desks failed: empty user id")
		return nil, core_errors.BadRequest()
	}

	desks, err := s.deskRepository.GetAllUsersDesks(ctx, userUUID)
	if err != nil {
		s.log.Error("repository get all users desks failed", slog.Any("userID", userUUID), slog.Any("err", err))
		return nil, err
	}

	s.log.Info("got all users desks successfully", slog.Any("userID", userUUID), slog.Int("desksCount", len(desks)))

	return desks, nil
}
