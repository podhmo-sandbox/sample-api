package repositorytest

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/model/entity"
	"github.com/podhmo/or"
)

type DBOption func(*sqlx.DB) error

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
		if err := opt(db); err != nil {
			t.Fatalf("error on %+v", err)
		}
	}
	return db, func() {
	}
}

func WithTodo(xs []entity.Todo) DBOption {
	// TODO: buildのところのファイルをもらってくる？
	return func(db *sqlx.DB) error {
		stmt := `
CREATE TABLE todo (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL
);
`
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("create todo table: %w", err)
		}

		// TODO: bulk insert
		for _, x := range xs {
			_, err := db.Exec("INSERT INTO todo (id, title, content) VALUES (?, ?, ?)", x.Id, x.Title, x.Content)
			if err != nil {
				return fmt.Errorf("insert data: %w", err)
			}
		}
		return nil
	}
}
