package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *DesksRepository) DeleteDesk(ctx context.Context, userUUID, deskUUID uuid.UUID) error {
	r.log.Info("deleting desk in repository", slog.Any("deskID", deskUUID))

	query := `DELETE FROM desks WHERE id = $1 AND owner_id = $2`

	if _, err := r.pool.Exec(ctx, query, deskUUID, userUUID); err != nil {
		r.log.Error("failed to delete desk in repository",
			slog.Any("deskID", deskUUID), slog.Any("userID", userUUID), slog.Any("err", err))
		return err
	}

	r.log.Info("desk deleted successfully in repository", slog.Any("deskID", deskUUID))

	return nil
}
