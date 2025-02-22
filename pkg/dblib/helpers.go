package dblib

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func AssertRowsCount(t *testing.T, db *sqlx.DB, tablename string, want int) {
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

func AssertRowsCountWith(t *testing.T, db *sqlx.DB, tablename string, want int) func(int) {
	t.Helper()
	AssertRowsCount(t, db, tablename, want)
	return func(want int) {
		t.Helper()
		AssertRowsCount(t, db, tablename, want)
	}
}
