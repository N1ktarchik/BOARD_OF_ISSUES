package repository

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *DesksRepository) ChangeDesksData(ctx context.Context, deskUpdate *domain.Desk, requesterID uuid.UUID) (*domain.Desk, error) {
	r.log.Info("starting desk update",
		slog.Any("deskID", deskUpdate.Id),
		slog.Any("requesterID", requesterID))

	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if deskUpdate.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argID))
		args = append(args, deskUpdate.Name)
		argID++
	}

	if deskUpdate.Password != "" {
		setValues = append(setValues, fmt.Sprintf("password = $%d", argID))
		args = append(args, deskUpdate.Password)
		argID++
	}

	if deskUpdate.OwnerId != uuid.Nil {
		setValues = append(setValues, fmt.Sprintf("owner_id = $%d", argID))
		args = append(args, deskUpdate.OwnerId)
		argID++
	}

	if len(setValues) == 0 {
		r.log.Warn("no fields to update in request", slog.Any("deskID", deskUpdate.Id))
		return nil, fmt.Errorf("no data to update")
	}

	query := fmt.Sprintf(
		"UPDATE desks SET %s WHERE id = $%d AND owner_id = $%d RETURNING id, name, password, owner_id, created_at",
		strings.Join(setValues, ", "),
		argID,
		argID+1,
	)

	args = append(args, deskUpdate.Id, requesterID)

	var updated deskModel
	row := r.pool.QueryRow(ctx, query, args...)

	if err := updated.scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn("update rejected: desk not found or unauthorized",
				slog.Any("deskID", deskUpdate.Id),
				slog.Any("requesterID", requesterID))

			return nil, core_errors.DeskNotFound()
		}

		r.log.Error("database error during desk update",
			slog.Any("deskID", deskUpdate.Id),
			slog.Any("err", err))

		return nil, core_errors.ServerError()
	}

	r.log.Info("desk updated successfully", slog.Any("deskID", updated.Id))

	domainDesk := modelToDomain(updated)
	return &domainDesk, nil
}
