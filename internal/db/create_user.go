package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) CreateUser(ctx context.Context, username, encryptedPass string) (*int, error) {
	insertQuery, insArgs, err := sq.Insert(usersTable).
		Columns(usersNameColumn, usersPasswordColumn).
		Values(username, encryptedPass).
		Suffix(fmt.Sprintf("RETURNING %s", userIDColumn)).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var userID *int
	row := s.db.QueryRowContext(ctx, insertQuery, insArgs...)
	err = row.Scan(&userID)
	if err != nil {
		return nil, err
	}

	return userID, nil
}
