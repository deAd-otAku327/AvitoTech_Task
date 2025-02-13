package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
)

var ErrNoUser = errors.New("no such user")

func (s *storage) SendCoinByUsername(ctx context.Context, userID int, destUsername string, amount int) error {
	updateSourceBalanceQuery, sourceUpdArgs, err := sq.Update(usersTable).
		Set(usersBalanceColumn, sq.Expr(fmt.Sprintf("%s - %d", usersBalanceColumn, amount))).
		Where(sq.Eq{userIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	updateDestBalanceQuery, destUpdArgs, err := sq.Update(usersTable).
		Set(usersBalanceColumn, sq.Expr(fmt.Sprintf("%s + %d", usersBalanceColumn, amount))).
		Where(sq.Eq{usersNameColumn: destUsername}).
		Suffix("RETURNING " + userIDColumn).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateSourceBalanceQuery, sourceUpdArgs...)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("tx rollback error: %s", txErr.Error())
		}
		return err
	}

	var destID int
	row := tx.QueryRowContext(ctx, updateDestBalanceQuery, destUpdArgs...)
	err = row.Scan(&destID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("tx rollback error: %s", txErr.Error())
		}
		if err == sql.ErrNoRows {
			return ErrNoUser
		}
		return err
	}

	// destID required.
	insertTransferQuery, insertArgs, err := sq.Insert(coinTransfersTable).
		Columns(coinTransfersSourceColumn, coinTransfersDestColumn, coinTransfersAmountColumn, coinTransfersTimeColumn).
		Values(userID, destID, amount, time.Now()).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("tx rollback error: %s", txErr.Error())
		}
		return err
	}

	_, err = tx.ExecContext(ctx, insertTransferQuery, insertArgs...)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("tx rollback error: %s", txErr.Error())
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
