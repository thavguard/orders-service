package db

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib" // Dont need an named import
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/otel/sdk/trace"
)

type DB struct {
	Pool *sqlx.DB
}

func NewDBConnection(ctx context.Context, tp *trace.TracerProvider, connString string) (*DB, error) {
	pool, err := otelsqlx.ConnectContext(ctx, "pgx", connString, otelsql.WithTracerProvider(tp))

	if err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() error {
	return db.Pool.Close()
}
