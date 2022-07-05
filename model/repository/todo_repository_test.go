package repository

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

func newDB(ctx context.Context, t *testing.T, options ...DBOption) (*sqlx.DB, func()) {
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

func withTodo(xs []entity.TodoEntity) DBOption {
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
		return nil
	}
}

func TestInsertTodo(t *testing.T) {
	ctx := context.Background()
	db, teardown := newDB(ctx, t, withTodo(nil))
	defer teardown()

	assertRowsCount(t, db, "todo", 0 /* want*/) // todo: checking by defer
	repo := &todoRepository{DB: db}
	id, err := repo.InsertTodo(entity.TodoEntity{Title: "go to bed", Content: "should sleep"})
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	_ = id
	assertRowsCount(t, db, "todo", 1 /* want*/)
}

func assertRowsCount(t *testing.T, db *sqlx.DB, tablename string, want int) {
	t.Helper()
	var got int
	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM "%s";`, tablename)
	if err := db.Get(&got, stmt); err != nil {
		t.Fatalf("assertRowsCount() unexpected error: %+v", err)
	}
	if want != got {
		t.Errorf("count(%v): want=%d != got=%d", tablename, want, got)
	}
}
