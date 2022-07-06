package dblib

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo/or"
)

type Config struct {
	Driver string `flag:"driver"`
	DSN    string `flag:"dsn"`
}

type DBOption func(*testing.T, *sqlx.DB)

func DefaultConfig() Config {
	return Config{Driver: "sqlite", DSN: ":memory:"}
}

func NewDB(ctx context.Context, t *testing.T, options ...DBOption) (*sqlx.DB, func()) {
	c := DefaultConfig()
	db := or.Fatal(sqlx.ConnectContext(ctx, c.Driver, c.DSN))(t)

	for _, opt := range options {
		opt(t, db)
	}
	return db, func() {
	}
}
