package db

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Pool *sqlx.DB
}

func InitDbConnection(ctx context.Context, connString string) (*DB, error) {
	pool, err := sqlx.Connect("pgx", connString)

	if err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
