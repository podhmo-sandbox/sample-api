package repository

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func init() {
	var err error
	// TODO: customizable
	Db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
}
