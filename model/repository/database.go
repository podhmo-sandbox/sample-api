package repository

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var Db *sqlx.DB

func init() {
	// TODO: customizable
	Db = sqlx.MustConnect("sqlite", ":memory:")
}
