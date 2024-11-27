package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"song-library/internal/constants"
)

type Database struct {
	*sql.DB
}

func NewDatabase(connStr string) (*Database, error) {
	db, err := sql.Open(constants.PostgresDriver, connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
