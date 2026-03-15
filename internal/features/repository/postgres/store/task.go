package store

import (
	"context"
	"time"

	repo "Board_of_issuses/internal/features/repository"
)

func (c *connect) CreateTask(ctx context.Context, task *repo.Task) error {
	query := `INSERT INTO tasks (userid,deskid,name,description,time) VALUES ($1,$2,$3,$4,$5) RETURNING created_at `

	err := c.db.QueryRow(ctx, query, task.UserId, task.DeskId, task.Name, task.Description, task.Time).Scan(&task.Created_at)
	if err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateTaskDecription(ctx context.Context, id int, description string) error {
	query := `UPDATE tasks SET description = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, description, id); err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateTaskTime(ctx context.Context, id int, time time.Time) error {
	query := `UPDATE tasks SET time = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, time, id); err != nil {
		return err
	}

	return nil
}

func (c *connect) UpdateTaskDone(ctx context.Context, id int) error {
	query := `UPDATE tasks SET done = $1 WHERE id = $2 `

	if _, err := c.db.Exec(ctx, query, true, id); err != nil {
		return err
	}

	return nil
}

func (c *connect) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1 `

	if _, err := c.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (c *connect) GetAllTasksFromOneDesk(ctx context.Context, deskId int) ([]repo.Task, error) {
	query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE deskid = $1`

	rows, err := c.db.Query(ctx, query, deskId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasksArr := make([]repo.Task, 0)

	for rows.Next() {
		task := repo.Task{}
		err := rows.Scan(
			&task.Id,
			&task.UserId,
			&task.DeskId,
			&task.Name,
			&task.Description,
			&task.Done,
			&task.Time,
			&task.Created_at,
		)

		if err != nil {
			return nil, err
		}

		tasksArr = append(tasksArr, task)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return tasksArr, nil
}

func (c *connect) GetTaskOwner(ctx context.Context, taskID int) (int, error) {
	query := `SELECT userid FROM tasks WHERE id =$1`

	var owner int

	if err := c.db.QueryRow(ctx, query, taskID).Scan(&owner); err != nil {
		return 0, err
	}

	return owner, nil
}

func (c *connect) GetDeskIDByTask(ctx context.Context, taskID int) (int, error) {
	query := `SELECT deskid FROM tasks WHERE id =$1`

	var deskid int

	if err := c.db.QueryRow(ctx, query, taskID).Scan(&deskid); err != nil {
		return 0, err
	}

	return deskid, nil
}

func (c *connect) GetTasksWithParams(ctx context.Context, deskId int, done bool) ([]repo.Task, error) {

	query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE deskid = $1 AND done =$2`

	rows, err := c.db.Query(ctx, query, deskId, done)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasksArr := make([]repo.Task, 0)

	for rows.Next() {
		task := repo.Task{}
		err := rows.Scan(
			&task.Id,
			&task.UserId,
			&task.DeskId,
			&task.Name,
			&task.Description,
			&task.Done,
			&task.Time,
			&task.Created_at,
		)

		if err != nil {
			return nil, err
		}

		tasksArr = append(tasksArr, task)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return tasksArr, nil
}
