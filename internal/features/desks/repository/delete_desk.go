package repository

import (
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *DesksRepository) DeleteDesk(ctx context.Context, userUUID, deskUUID uuid.UUID) error {
	r.log.Info("deleting desk in repository", slog.Any("deskID", deskUUID))

	query := `DELETE FROM desks WHERE id = $1 AND owner_id = $2`

	result, err := r.pool.Exec(ctx, query, deskUUID, userUUID)
	if err != nil {
		r.log.Error("failed to delete desk in repository",
			slog.Any("deskID", deskUUID), slog.Any("userID", userUUID), slog.Any("err", err))
		return core_errors.ServerError()
	}

	if result.RowsAffected() == 0 {
		r.log.Warn("user not owner of desk", slog.Any("userID", userUUID), slog.Any("deskID", deskUUID))

		return core_errors.UserNotOwnerOfDesk(userUUID.String(), deskUUID.String())
	}

	r.log.Info("desk deleted successfully in repository", slog.Any("deskID", deskUUID))

	return nil
}
