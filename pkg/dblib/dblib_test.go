package dblib

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func WithDummyTable() DBOption {
	return func(t *testing.T, db *sqlx.DB) {
		stmt := `
	CREATE TABLE dummy (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	);
	`
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("create dummy table: %v", err)
		}
	}
}

func TestCount(t *testing.T) {
	ctx := context.Background()
	db, teardown := NewDB(ctx, t, WithDummyTable())
	defer teardown()

	afterInsertCheck := AssertRowsCountWith(t, db, "dummy", 0 /* want */)
	defer afterInsertCheck(1 /* want */)

	if _, err := db.ExecContext(ctx, "INSERT INTO dummy(id,title,content) VALUES (?, ?, ?)", 1, "foo", ""); err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
}
