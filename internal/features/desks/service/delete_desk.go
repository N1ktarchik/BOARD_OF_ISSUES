package service

import (
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *DesksService) DeleteDesk(ctx context.Context, deskID, userID string) error {
	s.log.Info("deleting desk", slog.Any("deskID", deskID), slog.Any("userID", userID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.log.Warn("delete desk failed: error parsing user id", slog.Any("err", err))
		return core_errors.BadRequest()
	}
	if userUUID == uuid.Nil {
		s.log.Warn("delete desk failed: empty user id")
		return core_errors.BadRequest()
	}

	deskUUID, err := uuid.Parse(deskID)
	if err != nil {
		s.log.Warn("delete desk failed: error parsing desk id", slog.Any("userID", userUUID), slog.Any("err", err))
		return core_errors.BadRequest()
	}
	if deskUUID == uuid.Nil {
		s.log.Warn("delete desk failed: empty desk id", slog.Any("userID", userUUID))
		return core_errors.BadRequest()
	}

	if err := s.deskRepository.DeleteDesk(ctx, userUUID, deskUUID); err != nil {
		s.log.Error("repository delete desk failed",
			slog.Any("deskID", deskUUID), slog.Any("userID", userUUID), slog.Any("err", err))

		return err
	}

	s.log.Info("desk deleted successfully", slog.Any("deskID", deskUUID), slog.Any("userID", userUUID))
	return nil
}
