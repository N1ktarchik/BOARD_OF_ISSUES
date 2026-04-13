package postgres

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDB(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv("CONSTR")
	if connStr == "" {
		return nil, errors.New("CONSTR is not set")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err := createTables(ctx, pool); err != nil {
		return nil, err
	}

	return pool, nil

}

func createTables(ctx context.Context, db *pgxpool.Pool) error {
	tables := []string{

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			login VARCHAR(100) NOT NULL UNIQUE,
			password TEXT NOT NULL, -- Используем TEXT для хэша
			email VARCHAR(255) UNIQUE, -- UNIQUE для почты
			name VARCHAR(100) NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc')
		);`,

		`CREATE TABLE IF NOT EXISTS desks (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL DEFAULT 'new desk',
			password TEXT NOT NULL, -- Код доступа к доске
			owner_id UUID NOT NULL,
			created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc'),

			CONSTRAINT fk_owner FOREIGN KEY (owner_id) 
				REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS desk_members (
			user_id UUID NOT NULL,
			desk_id UUID NOT NULL,
			
			PRIMARY KEY (user_id, desk_id),

			CONSTRAINT fk_member_user FOREIGN KEY (user_id) 
				REFERENCES users(id) ON DELETE CASCADE,
			
			CONSTRAINT fk_member_desk FOREIGN KEY (desk_id) 
				REFERENCES desks(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY,
			author_id UUID NOT NULL,
			desk_id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT DEFAULT '',
			done BOOLEAN NOT NULL DEFAULT FALSE,
			deadline TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc'),

			CONSTRAINT fk_task_desk FOREIGN KEY (desk_id) 
				REFERENCES desks(id) ON DELETE CASCADE,

			CONSTRAINT fk_task_author FOREIGN KEY (author_id) 
				REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS users_idx ON users(id);`,
		`CREATE INDEX IF NOT EXISTS deskusers_idx ON desksusers(userid,deskid);`,
		`CREATE INDEX IF NOT EXISTS desks_idx ON desks(id);`,
		`CREATE INDEX IF NOT EXISTS tasks_idx ON tasks(id);`,
		`CREATE INDEX IF NOT EXISTS tasks_help_idx ON tasks(userid,deskid);`,
	}

	for _, query := range tables {
		if _, err := db.Exec(ctx, query); err != nil {
			return err
		}
	}

	for _, query := range indexes {
		if _, err := db.Exec(ctx, query); err != nil {
			return err
		}
	}

	return nil
}
