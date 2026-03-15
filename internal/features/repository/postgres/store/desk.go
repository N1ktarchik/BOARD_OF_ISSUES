package store

import (
	"context"
	"database/sql"

	repo "Board_of_issuses/internal/features/repository"
)

func (c *connect) CreateDesk(ctx context.Context, desk *repo.Desk) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `INSERT INTO desks (name,password,ownerid) VALUES ($1,$2,$3) RETURNING id`
	err = tx.QueryRow(ctx, query, desk.Name, desk.Password, desk.OwnerId).Scan(&desk.Id)
	if err != nil {
		return err
	}

	query = `INSERT INTO desksusers (userid,deskid) VALUES ($1,$2)`
	result, err := tx.Exec(ctx, query, desk.OwnerId, desk.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	tx.Commit(ctx)
	return nil
}

func (c *connect) UpdateDeskName(ctx context.Context, deskId int, name string) error {
	query := `UPDATE desks SET name = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, name, deskId); err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateDesksPassword(ctx context.Context, deskId int, password string) error {
	query := `UPDATE desks SET password = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, password, deskId); err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateDeskOwner(ctx context.Context, ownerid, deskid int) error {
	query := `UPDATE desks SET ownerid = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, ownerid, deskid); err != nil {
		return err
	}

	return nil
}

func (c *connect) DeleteDesk(ctx context.Context, deskId int) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM desksusers WHERE deskid = $1`
	result, err := tx.Exec(ctx, query, deskId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	query = `DELETE FROM desks WHERE id = $1`
	result, err = c.db.Exec(ctx, query, deskId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	tx.Commit(ctx)

	return nil
}

func (c *connect) CheckDeskOwner(ctx context.Context, deskId int) (int, error) {
	query := `SELECT ownerid FROM desks WHERE id = $1 `
	var ownerID int

	if err := c.db.QueryRow(ctx, query, deskId).Scan(&ownerID); err != nil {
		return 0, nil
	}

	return ownerID, nil
}

func (c *connect) CheckDeskPassword(ctx context.Context, deskId int) (string, error) {
	query := `SELECT password FROM desks WHERE id = $1 `
	var password string

	if err := c.db.QueryRow(ctx, query, deskId).Scan(&password); err != nil {
		return "", nil
	}

	return password, nil
}
