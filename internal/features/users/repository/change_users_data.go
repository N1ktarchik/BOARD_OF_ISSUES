package repository

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *UsersRepository) ChangeUsersData(ctx context.Context, user *domain.User) (*domain.User, error) {

	query := "UPDATE users SET "
	var args []any
	argId := 1

	if user.Name != "" {
		query += fmt.Sprintf("name=$%d, ", argId)
		args = append(args, user.Name)
		argId++
	}
	if user.Email != "" {
		query += fmt.Sprintf("email=$%d, ", argId)
		args = append(args, user.Email)
		argId++
	}
	if user.Password != "" {
		query += fmt.Sprintf("password=$%d, ", argId)
		args = append(args, user.Password)
		argId++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id=$%d", argId)
	args = append(args, user.ID)
	query += " RETURNING id, login, password, email, name, created_at"


	var model userModel
	row := r.pool.QueryRow(ctx, query, args...)

	if err := model.scan(row); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {

				r.log.Error("change users data failed : email already registered",
					slog.String("email", user.Email))

				return nil, core_errors.EmailAlreadyRegistered(user.Email)
			}
		}

		r.log.Error("change users data failed : database error", slog.Any("err", err))
		return nil, core_errors.ServerError()
	}

	result := modelToDomain(model)
	r.log.Info("users data changed successfully", slog.String("user_id", user.ID.String()))

	return &result, nil
}
