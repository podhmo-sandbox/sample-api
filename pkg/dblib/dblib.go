package dblib

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo/or"
)

type DBOption func(*testing.T, *sqlx.DB)

type DBConfig struct {
	Driver string
	DSN    string
}

func DefaultDBConfig() DBConfig {
	return DBConfig{Driver: "sqlite", DSN: ":memory:"}
}

func NewDB(ctx context.Context, t *testing.T, options ...DBOption) (*sqlx.DB, func()) {
	c := DefaultDBConfig()
	db := or.Fatal(sqlx.ConnectContext(ctx, c.Driver, c.DSN))(t)

	for _, opt := range options {
		opt(t, db)
	}
	return db, func() {
	}
}
