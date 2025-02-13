package db

import (
	"context"
	"database/sql"
	"merch_shop/internal/config"

	_ "github.com/lib/pq"
)

type DB interface {
	CreateOrGetUser(ctx context.Context, username, encryptedPass string) (*int, string, error)
}

const (
	usersTable          = "users"
	userIDColumn        = "id"
	usersNameColumn     = "username"
	usersPasswordColumn = "password"
)

type storage struct {
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

	database.SetMaxOpenConns(cfg.MaxOpenConns)

	return &storage{db: database}, nil
}
