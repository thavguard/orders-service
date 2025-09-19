package db

import (
	_ "github.com/jackc/pgx/v5/stdlib" // Dont need an named import
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Pool *sqlx.DB
}

func NewDBConnection(connString string) (*DB, error) {
	pool, err := sqlx.Connect("pgx", connString)

	if err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() error {
	return db.Pool.Close()
}
