// Package sql provides a generic interface around SQL objects.
package sql

import (
	"database/sql"

	_ "github.com/lib/pq" // need for postgres driver
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

// NewDBFromConnStr creates a new data connection handle from a given
// connection string.
func NewDBFromConnStr(connStr *string) (*reform.DB, error) {
	conn, err := sql.Open("postgres", *connStr)
	if err == nil {
		err = conn.Ping()
	}
	if err != nil {
		return nil, err
	}

	dummy := func(format string, args ...interface{}) {}

	return reform.NewDB(conn,
		postgresql.Dialect, reform.NewPrintfLogger(dummy)), nil
}

// CloseDB closes database connection.
func CloseDB(db *reform.DB) {
	db.DBInterface().(*sql.DB).Close()
}
