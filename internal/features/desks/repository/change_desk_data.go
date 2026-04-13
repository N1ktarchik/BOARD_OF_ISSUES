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
	r.log.Info("updating desk data in repository", slog.Any("deskID", deskUpdate.Id), slog.Any("requesterID", requesterID))

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
		return nil, fmt.Errorf("no data to update")
	}

	query := fmt.Sprintf(
		"UPDATE desks SET %s WHERE id = $%d AND owner_id = $%d RETURNING id, name, password, owner_id, created_at",
		strings.Join(setValues, ", "),
		argID,
		argID+1,
	)

	args = append(args, deskUpdate.Id, deskUpdate.OwnerId)

	var updated deskModel
	row := r.pool.QueryRow(ctx, query, args...)

	if err := updated.scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core_errors.DeskNotFound()
		}
		return nil, core_errors.ServerError()
	}

	domainDesk := modelToDomain(updated)
	return &domainDesk, nil
}
