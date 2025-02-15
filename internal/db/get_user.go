package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) GetUser(ctx context.Context, username string) (*int, string, error) {
	selectQuery, selArgs, err := sq.Select(userIDColumn, usersPasswordColumn).
		From(usersTable).
		Where(sq.Eq{usersNameColumn: username}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, "", err
	}

	var userID *int
	var dbPassword string
	row := s.db.QueryRowContext(ctx, selectQuery, selArgs...)
	err = row.Scan(&userID, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", ErrNoUser
		}
		return nil, "", err
	}

	return userID, dbPassword, nil
}
