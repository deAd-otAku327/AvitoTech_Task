package db

import (
	"database/sql"
	"merch_shop/internal/config"

	_ "github.com/lib/pq"
)

type DB interface {
}

type Storage struct {
	db *sql.DB
}

func New(cfg config.DB) (DB, error) {
	database, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}
