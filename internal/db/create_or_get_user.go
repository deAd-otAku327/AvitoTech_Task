package db

import (
	"context"
	"database/sql"
	"fmt"
	"merch_shop/pkg/cryptor"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) CreateOrGetUser(ctx context.Context, username, password string) (*int, string, error) {
	selectQuery, selArgs, err := sq.Select(userIDColumn, usersPasswordColumn).
		From(usersTable).
		Where(sq.Eq{usersNameColumn: username}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, "", err
	}

	var userID *int
	var dbPassword *string
	row := s.db.QueryRowContext(ctx, selectQuery, selArgs...)
	err = row.Scan(&userID, &dbPassword)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", err
	}

	if err == sql.ErrNoRows {
		encryptedPass, err := cryptor.EncryptKeyword(password)
		if err != nil {
			return nil, "", err
		}

		insertQuery, insArgs, err := sq.Insert(usersTable).
			Columns(usersNameColumn, usersPasswordColumn).
			Values(username, encryptedPass).
			Suffix(fmt.Sprintf("RETURNING %s, %s", userIDColumn, usersPasswordColumn)).
			PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, "", err
		}

		row = s.db.QueryRowContext(ctx, insertQuery, insArgs...)
		err = row.Scan(&userID, &dbPassword)
		if err != nil {
			return nil, "", err
		}
	}

	return userID, *dbPassword, nil
}
