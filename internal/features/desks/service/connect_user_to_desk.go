package service

import (
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *DesksService) ConnectUserToDesk(ctx context.Context, userID, deskID uuid.UUID) error {
	s.log.Info("connecting user to desk", slog.Any("deskID", deskID), slog.Any("userID", userID))

	if userID == uuid.Nil {
		s.log.Error("connect user to desk failed: empty user id")
		return core_errors.BadRequest()
	}

	if deskID == uuid.Nil {
		s.log.Error("connect user to desk failed: empty desk id")
		return core_errors.BadRequest()
	}
	
	if err := s.deskRepository.ConnectUserToDesk(ctx, userID, deskID); err != nil {
		s.log.Error("repository connect user to desk failed", slog.Any("deskID", deskID),
			slog.Any("userID", userID), slog.Any("err", err))

		return err
	}

	s.log.Info("user connected to desk successfully", slog.Any("deskID", deskID), slog.Any("userID", userID))
	return nil

}
