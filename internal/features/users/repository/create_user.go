package repository

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (id,login, password, email, name) 
	VALUES ($1, $2, $3, $4,$5) `

	if _, err := r.pool.Exec(ctx, query, user.ID, user.Login, user.Password, user.Email, user.Name); err != nil {
		r.log.Error("failed to create user in database", slog.String("email", user.Email), slog.Any("err", err))

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			r.log.Error("user already registered", slog.String("email", user.Email), slog.String("login", user.Login))
			return core_errors.UserAlreadyRegistered(user.Login, user.Email)
		}

		return core_errors.ServerError()
	}

	r.log.Info("user successfully created in database", slog.String("user_id", user.ID.String()),
		slog.String("login", user.Login),
	)

	return nil
}
