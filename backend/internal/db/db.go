package db

import "database/sql"

var DB *sql.DB

type DatabaseLike interface {
	Exec(query string, args ...any) (sql.Result, error)
}
