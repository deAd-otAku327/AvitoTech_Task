package db

import (
	"context"
	"database/sql"
	"errors"
	"merch_shop/internal/config"
	"merch_shop/internal/models"

	_ "github.com/lib/pq"
)

type DB interface {
	CreateOrGetUser(ctx context.Context, username, encryptedPass string) (*int, string, error)
	SendCoinByUsername(ctx context.Context, userID int, destUsername string, amount int) error
	BuyItemByItemID(ctx context.Context, userID, itemID int) error
	GetUserInfoByUserID(ctx context.Context, userID int) (*int, []byte, *models.CoinTransferHistory, error)
}

const (
	usersTable           = "users"
	userIDColumn         = "id"
	usersNameColumn      = "username"
	usersPasswordColumn  = "password"
	usersBalanceColumn   = "balance"
	usersInventoryColumn = "inventory"

	coinTransfersTable        = "coin_transfers"
	coinTransfersSourceColumn = "from_user_id"
	coinTransfersDestColumn   = "to_user_id"
	coinTransfersAmountColumn = "amount"
	coinTransfersTimeColumn   = "timing"

	itemsTable       = "items"
	itemsIDColumn    = "id"
	itemsTypeColumn  = "type"
	itemsPriceColumn = "price"
)

var (
	ErrNoUser         = errors.New("no such user")
	ErrNoItem         = errors.New("no such item")
	ErrNotEnoughCoins = errors.New("not enough coins")
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
