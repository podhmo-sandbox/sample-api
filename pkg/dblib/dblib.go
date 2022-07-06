package dblib

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo/or"
)

type Config struct {
	Driver string `flag:"driver"`
	DSN    string `flag:"dsn"`
}

func (c *Config) New(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, c.Driver, c.DSN)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	return db, nil
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
