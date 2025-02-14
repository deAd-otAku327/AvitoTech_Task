package db

import (
	"context"
	"fmt"
	"merch_shop/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) GetUserInfoByUserID(ctx context.Context, userID int) (*int, []byte, *models.CoinTransferHistory, error) {
	selectUserDataQuery, userSelectArgs, err := sq.Select(usersBalanceColumn, usersInventoryColumn).
		From(usersTable).
		Where(sq.Eq{userIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, nil, nil, err
	}

	selectOutgoingTransfersQuery, outgoingTransfersArgs, err := sq.Select(usersNameColumn, coinTransfersAmountColumn).
		From(coinTransfersTable).
		LeftJoin(fmt.Sprintf("%s ON %s.%s = %s.%s", usersTable, coinTransfersTable, coinTransfersDestColumn, usersTable, userIDColumn)).
		Where(sq.Eq{coinTransfersSourceColumn: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, nil, nil, err
	}

	selectIngoingTransfersQuery, ingoingTransfersArgs, err := sq.Select(usersNameColumn, coinTransfersAmountColumn).
		From(coinTransfersTable).
		LeftJoin(fmt.Sprintf("%s ON %s.%s = %s.%s", usersTable, coinTransfersTable, coinTransfersSourceColumn, usersTable, userIDColumn)).
		Where(sq.Eq{coinTransfersDestColumn: userID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, nil, nil, err
	}

	var balance *int
	var inventory []byte

	row := s.db.QueryRowContext(ctx, selectUserDataQuery, userSelectArgs...)
	err = row.Scan(&balance, &inventory)
	if err != nil {
		return nil, nil, nil, err
	}

	history := &models.CoinTransferHistory{
		Recieved: make([]models.IngoingCoinTransfer, 0),
		Sent:     make([]models.OutgoingCoinTransfer, 0),
	}
	inTransfer := models.IngoingCoinTransfer{}
	outTransfer := models.OutgoingCoinTransfer{}

	rowsOut, err := s.db.QueryContext(ctx, selectOutgoingTransfersQuery, outgoingTransfersArgs...)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rowsOut.Close()

	for rowsOut.Next() {
		err := rowsOut.Scan(&outTransfer.Username, &outTransfer.Amount)
		if err != nil {
			return nil, nil, nil, err
		}
		history.Sent = append(history.Sent, outTransfer)
	}
	if rowsOut.Err() != nil {
		return nil, nil, nil, err
	}

	rowsIn, err := s.db.QueryContext(ctx, selectIngoingTransfersQuery, ingoingTransfersArgs...)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rowsIn.Close()

	for rowsIn.Next() {
		err := rowsIn.Scan(&inTransfer.Username, &inTransfer.Amount)
		if err != nil {
			return nil, nil, nil, err
		}
		history.Recieved = append(history.Recieved, inTransfer)
	}
	if rowsIn.Err() != nil {
		return nil, nil, nil, err
	}

	return balance, inventory, history, nil
}
