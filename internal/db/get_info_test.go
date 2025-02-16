package db

import (
	"context"
	"log"
	"merch_shop/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	selectUserDataQueryRegexp = `
		SELECT (.*) FROM users WHERE (.*)
	`
	selectTransfersQueryRegexp = `
		SELECT (.*) FROM coin_transfers LEFT JOIN users ON (.*) WHERE (.*)
	`
)

func TestGetUserInfoByUserID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer mockDB.Close()

	db := storage{db: mockDB}

	type mockBehavior func(args int)

	type expectedRes struct {
		balance   *int
		inventory []byte
		history   *models.CoinTransferHistory
	}

	expBalance := 1000

	testCases := []struct {
		name       string
		expected   expectedRes
		dbBehavior mockBehavior

		expectErr bool
	}{
		{
			name: "positive result",
			expected: expectedRes{
				balance:   &expBalance,
				inventory: []byte("test"),
				history: &models.CoinTransferHistory{
					Recieved: []models.IngoingCoinTransfer{{
						Username: "testIn",
						Amount:   200,
					}},
					Sent: []models.OutgoingCoinTransfer{{
						Username: "testOut",
						Amount:   100,
					}},
				},
			},
			dbBehavior: func(arg int) {
				selectUserDataRows := sqlmock.NewRows([]string{usersBalanceColumn, usersInventoryColumn}).AddRow(expBalance, []byte("test"))
				mock.ExpectQuery(selectUserDataQueryRegexp).WithArgs(arg).WillReturnRows(selectUserDataRows)

				selectOutgoingTransfersRows := sqlmock.NewRows([]string{usersBalanceColumn, coinTransfersAmountColumn}).AddRow(
					"testOut", 100,
				)
				mock.ExpectQuery(selectTransfersQueryRegexp).WithArgs(arg).WillReturnRows(selectOutgoingTransfersRows)

				selectIngoingTransfersRows := sqlmock.NewRows([]string{usersBalanceColumn, coinTransfersAmountColumn}).AddRow(
					"testIn", 200,
				)
				mock.ExpectQuery(selectTransfersQueryRegexp).WithArgs(arg).WillReturnRows(selectIngoingTransfersRows)
			},
			expectErr: false,
		},
	}

	userID := 0

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.dbBehavior(userID)

			balance, inventory, history, err := db.GetUserInfoByUserID(context.Background(), userID)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tc.expected.balance, balance)
				assert.Equal(t, tc.expected.inventory, inventory)
				assert.Equal(t, tc.expected.history, history)
			}
		})
	}

}
