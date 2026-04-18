package repository

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *DesksRepository) GetAllUsersDesks(ctx context.Context, userUUID uuid.UUID) ([]domain.Desk, error) {
	r.log.Info("fetching all desks for user", slog.Any("userID", userUUID))

	query := `
        SELECT d.id, d.name, d.password, d.owner_id, d.created_at 
        FROM desks d
        JOIN desks_users du ON d.id = du.desk_id
        WHERE du.user_id = $1
        ORDER BY d.created_at DESC`

	rows, err := r.pool.Query(ctx, query, userUUID)
	if err != nil {
		r.log.Error("failed to query desks", slog.Any("userID", userUUID), slog.Any("err", err))
		return nil, core_errors.ServerError()
	}
	defer rows.Close()

	desks := make([]domain.Desk, 0)

	for rows.Next() {
		var d deskModel

		err := rows.Scan(&d.Id, &d.Name, &d.Password, &d.OwnerId, &d.Created_at)
		if err != nil {
			r.log.Error("failed to scan desk row", slog.Any("err", err))
			return nil, core_errors.ServerError()
		}

		desks = append(desks, modelToDomain(d))
	}

	if err = rows.Err(); err != nil {
		return nil, core_errors.ServerError()
	}

	r.log.Info("successfully fetched desks", slog.Int("count", len(desks)))
	return desks, nil
}
