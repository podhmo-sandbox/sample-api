package repositorytest

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/entity"
	"github.com/podhmo-sandbox/sample-api/pkg/dblib"
)

// setup

type DBOption = dblib.DBOption

var (
	DefaultDBConfig = dblib.DefaultDBConfig
	NewDB           = dblib.NewDB
)

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
