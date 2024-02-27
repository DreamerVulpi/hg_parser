package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectToPostgres(driver, params string) (*sql.DB, error) {
	conn, err := sql.Open(driver, params)
	if err != nil {
		return nil, err
	}
	return conn, err
}
