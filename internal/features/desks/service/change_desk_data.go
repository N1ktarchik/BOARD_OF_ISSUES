package service

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *DesksService) ChangeDesksData(ctx context.Context, deskUpdate *domain.Desk, requesterID uuid.UUID) (*domain.Desk, error) {
	s.log.Info("changing desk data", slog.Any("deskID", deskUpdate.Id), slog.Any("userID", deskUpdate.OwnerId))

	if deskUpdate.Name != "" && len(deskUpdate.Name) < 3 {
		s.log.Error("change desk data failed: short desk name", slog.Any("deskName", deskUpdate.Name),
			slog.Any("deskID", deskUpdate.Id), slog.Any("userID", deskUpdate.OwnerId))

		return nil, core_errors.BadRequest()
	}

	if requesterID == uuid.Nil {
		s.log.Error("change desk data failed: empty requester id", slog.Any("deskID", deskUpdate.Id),
			slog.Any("userID", deskUpdate.OwnerId))

		return nil, core_errors.BadRequest()
	}

	if deskUpdate.Password != "" {
		hashPassword, err := domain.Hash(deskUpdate.Password)
		if err != nil {
			s.log.Error("change desk data failed: password hashing error", slog.Any("err", err))
			return nil, err
		}

		deskUpdate.Password = hashPassword
	}

	saveDesk, err := s.deskRepository.ChangeDesksData(ctx, deskUpdate, requesterID)
	if err != nil {
		s.log.Error("change desk data failed: repository error", slog.Any("deskID", deskUpdate.Id),
			slog.Any("userID", deskUpdate.OwnerId), slog.Any("err", err))

		return nil, err
	}

	return saveDesk, nil

}
