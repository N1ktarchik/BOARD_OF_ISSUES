package repository

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUser(ctx context.Context, email, login string) (*domain.User, error) {
	query := `SELECT id,login,password,email,name,created_at FROM users WHERE login = $1 OR email = $2`

	row := r.pool.QueryRow(ctx, query, login, email)

	userModel := userModel{}
	if err := userModel.scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Error("user not found", slog.String("email", email), slog.String("login", login))
			return nil, core_errors.UserNotFound()
		}

		r.log.Error("failed to get user", slog.String("email", email), slog.String("login", login), slog.Any("err", err))
		return nil, core_errors.ServerError()
	}

	saveUser := modelToDomain(userModel)

	return &saveUser, nil

}
