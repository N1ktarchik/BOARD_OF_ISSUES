package store

import (
	"context"
	"database/sql"
)

func (c *connect) ConnectUserToDesk(ctx context.Context, userID, deskID int) error {
	query := `INSERT INTO desksusers (userid,deskid) VALUES ($1,$2)`

	if _, err := c.db.Exec(ctx, query, userID, deskID); err != nil {
		return err
	}

	return nil
}

func (c *connect) GetUserDesks(ctx context.Context, userId int) ([]int, error) {
	query := `SELECT deskid FROM desksusers WHERE userid = $1`

	deskArr := make([]int, 0)

	rows, err := c.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var deskId int

		if err := rows.Scan(&deskId); err != nil {
			return nil, err
		}

		deskArr = append(deskArr, deskId)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return deskArr, nil
}

func (c *connect) DeleteUserDesk(ctx context.Context, userId, deskId int) error {
	query := `DELETE FROM desksusers WHERE userid = $1 AND deskid = $2`

	result, err := c.db.Exec(ctx, query, userId, deskId)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (c *connect) CheckUserDesk(ctx context.Context, userId, deskId int) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM desksusers WHERE userid = $1 AND deskid = $2)`

	var exists bool

	err := c.db.QueryRow(ctx, query, userId, deskId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
