package repository

import (
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *DesksRepository) ConnectUserToDesk(ctx context.Context, userID, deskID uuid.UUID) error {
	r.log.Info("connecting user to desk in repository", slog.Any("deskID", deskID), slog.Any("userID", userID))

	query := `INSERT INTO desk_members (user_id, desk_id) VALUES ($1, $2)`

	if _, err := r.pool.Exec(ctx, query, userID, deskID); err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgForeignKeyViolation {
				r.log.Error("failed to connect user to desk: desk not found", slog.Any("err", err))

				return core_errors.DeskNotFound()
			}
		}

		r.log.Error("failed to connect user to desk: server error",
			slog.Any("userID", userID), slog.Any("deskID", deskID), slog.Any("err", err))

		return core_errors.ServerError()
	}

	r.log.Info("user connected to desk successfully in repository", slog.Any("deskID", deskID), slog.Any("userID", userID))

	return nil

}
