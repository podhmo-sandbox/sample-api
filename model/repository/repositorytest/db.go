package repositorytest

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/model/entity"
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

func WithTodo(xs []entity.Todo) DBOption {
	// TODO: buildのところのファイルをもらってくる？
	return func(t *testing.T, db *sqlx.DB) {
		stmt := `
CREATE TABLE todo (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL
);
`
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("create todo table: %v", err)
		}

		// TODO: bulk insert
		if len(xs) > 0 {
			t.Logf("\tinsert %d rows.", len(xs))
			for i, x := range xs {
				_, err := db.Exec("INSERT INTO todo (id, title, content) VALUES (?, ?, ?)", x.ID, x.Title, x.Content)
				if err != nil {
					t.Fatalf("insert data(%d): %+v", i, err)
				}
			}
		}
	}
}
