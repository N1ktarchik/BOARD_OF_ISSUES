package store

import (
	"context"
	"database/sql"

	repo "Board_of_issuses/internal/features/repository"
)

func (c *connect) CreateUser(ctx context.Context, user *repo.User) (int, error) {
	query := `INSERT INTO users (login, password, email, name) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int

	if err := c.db.QueryRow(ctx, query, user.Login, user.Password, user.Email, user.Name).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *connect) UpdateUserEmail(ctx context.Context, email string, userId int) error {
	query := `UPDATE users SET email = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, email, userId); err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateUserName(ctx context.Context, name string, userId int) error {
	query := `UPDATE users SET name = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, name, userId); err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateUserPassword(ctx context.Context, password string, userId int) error {
	query := `UPDATE users SET password = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, password, userId); err != nil {
		return err
	}

	return nil
}

func (c *connect) DeleteUser(ctx context.Context, userId string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `DELETE FROM desksusers WHERE userid = $1`
	result, err := tx.Exec(ctx, query, userId)
	if err != nil {
		return err
	}

	query = `DELETE FROM users WHERE id = $1`
	result, err = tx.Exec(ctx, query, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	tx.Commit(ctx)

	return nil
}

func (c *connect) GetUserByID(ctx context.Context, id int) (*repo.User, error) {
	query := `SELECT id,login,password,email,name,created_at FROM users WHERE id = $1`

	user := &repo.User{}

	err := c.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Email,
		&user.Name,
		&user.Created_at,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (c *connect) GetUserByLoginOrEmail(ctx context.Context, login, email string) (*repo.User, error) {
	query := `SELECT id,login,password,email,name,created_at FROM users WHERE login = $1 OR email = $2`

	user := &repo.User{}

	err := c.db.QueryRow(ctx, query, login, email).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Email,
		&user.Name,
		&user.Created_at,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (c *connect) CheckUserByEmailOrLogin(ctx context.Context, login, email string) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM users WHERE login = $1 OR email = $2)`

	var exists bool

	err := c.db.QueryRow(ctx, query, login, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
