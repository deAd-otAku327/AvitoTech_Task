package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) CreateOrGetUser(ctx context.Context, username, encryptedPass string) (*int, string, error) {
	insertQuery, insArgs, err := sq.Insert(usersTable).
		Columns(usersNameColumn, usersPasswordColumn).
		Values(username, encryptedPass).
		Suffix("ON CONFLICT (username) DO NOTHING").
		Suffix("RETURNING id, password").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, "", err
	}

	selectQuery, selArgs, err := sq.Select(userIDColumn, usersPasswordColumn).
		From(usersTable).
		Where(sq.Eq{usersNameColumn: username}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, "", err
	}

	var userID *int
	var dbPassword *string

	row := s.db.QueryRowContext(ctx, insertQuery, insArgs...)

	err = row.Scan(&userID, &dbPassword)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", err
	}

	if err == sql.ErrNoRows {
		row := s.db.QueryRowContext(ctx, selectQuery, selArgs...)

		err := row.Scan(&userID, &dbPassword)

		if err != nil {
			return nil, "", err
		}
	}

	return userID, *dbPassword, nil
}
