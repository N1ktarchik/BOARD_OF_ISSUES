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

		`CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			login VARCHAR(200) NOT NULL,
			password VARCHAR(200) NOT NULL,
			email VARCHAR(200) DEFAULT '',
			name VARCHAR(100) NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

		
		)`,

		`CREATE TABLE IF NOT EXISTS desksusers(
			userid SERIAL NOT NULL,
			deskid SERIAL NOT NULL
		
		)`,

		`CREATE TABLE IF NOT EXISTS desks(
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL DEFAULT 'userdesk',
				password VARCHAR(100) NOT NULL,
				ownerid SERIAL NOT NULL,
				created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		
		)`,

		`CREATE TABLE IF NOT EXISTS tasks(
				id SERIAL PRIMARY KEY,
				userid SERIAL NOT NULL,
				deskid SERIAL NOT NULL,
				name VARCHAR(100) NOT NULL,
				description VARCHAR(255) DEFAULT '',
				done BOOLEAN NOT NULL DEFAULT FALSE,
				time TIMESTAMP NOT NULL,
				created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)`,
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
