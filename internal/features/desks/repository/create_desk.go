package repository

import (
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"log/slog"

	"N1ktarchik/Board_of_issues/internal/core/domain"
	"context"
)

func (r *DesksRepository) CreateDesk(ctx context.Context, desk *domain.Desk) (*domain.Desk, error) {
	r.log.Info("creating desk", slog.Any("deskID", desk.Id), slog.Any("userID", desk.OwnerId))

	query1 := `INSERT INTO desks (id, name, password, owner_id) VALUES ($1, $2, $3, $4)
	 RETURNING id, name, password, owner_id, created_at`

	query2 := `INSERT INTO desk_members (user_id, desk_id) VALUES ($1, $2)`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error("failed to begin transaction", slog.Any("err", err))
		return nil, core_errors.ServerError()
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	row := tx.QueryRow(ctx, query1, desk.Id, desk.Name, desk.Password, desk.OwnerId)
	createdDesk := deskModel{}

	if err := createdDesk.scan(row); err != nil {
		r.log.Error("failed to create desk", slog.Any("deskID", desk.Id), slog.Any("userID", desk.OwnerId), slog.Any("err", err))
		return nil, core_errors.ServerError()
	}

	if _, err := tx.Exec(ctx, query2, desk.OwnerId, createdDesk.Id); err != nil {
		r.log.Error("failed to associate user with desk",
			slog.Any("userID", desk.OwnerId), slog.Any("deskID", createdDesk.Id), slog.Any("err", err))

		return nil, core_errors.ServerError()
	}

	domainDesk := modelToDomain(createdDesk)

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error("failed to commit transaction", slog.Any("err", err))

		return nil, core_errors.ServerError()
	}

	r.log.Info("desk created successfully", slog.Any("deskID", domainDesk.Id), slog.Any("userID", domainDesk.OwnerId))

	return &domainDesk, nil
}
