package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func (s *storage) BuyItemByItemID(ctx context.Context, userID, itemID int) error {

	selectItemQuery, itemSelectArgs, err := sq.Select(itemsTypeColumn, itemsPriceColumn).
		From(itemsTable).
		Where(sq.Eq{itemsIDColumn: itemID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	var itemType string
	var itemPrice int
	row := s.db.QueryRowContext(ctx, selectItemQuery, itemSelectArgs...)
	err = row.Scan(&itemType, &itemPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoItem
		}
		return err
	}

	updateUserBalanceQuery, balanceUpdArgs, err := sq.Update(usersTable).
		Set(usersBalanceColumn, sq.Expr(
			fmt.Sprintf("%s - %d", usersBalanceColumn, itemPrice))).
		Where(sq.Eq{userIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	updateUserInventoryQuery, inventoryUpdArgs, err := sq.Update(usersTable).
		Set(usersInventoryColumn, sq.Expr(
			fmt.Sprintf("jsonb_set(%s, '{%s}', (SELECT (COALESCE(%s -> '%s', '0'))::int + 1 FROM %s WHERE %s = %d)::text::jsonb)",
				usersInventoryColumn, itemType, usersInventoryColumn, itemType, usersTable, userIDColumn, userID))).
		Where(sq.Eq{userIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateUserBalanceQuery, balanceUpdArgs...)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("tx rollback error: %s", txErr.Error())
		}
		if pqErr, ok := err.(*pq.Error); ok {
			// From http://www.postgresql.org/docs/9.3/static/errcodes-appendix.html
			if pqErr.Code.Name() == "check_violation" {
				return ErrNotEnoughCoins
			}
		}
		return err
	}

	_, err = tx.ExecContext(ctx, updateUserInventoryQuery, inventoryUpdArgs...)
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
