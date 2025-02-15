package service

import (
	"context"
	"errors"
	"log/slog"
	"merch_shop/internal/db"
	dbmock "merch_shop/internal/db/mocks"
	"merch_shop/internal/models"
	cryptormock "merch_shop/pkg/cryptor/mocks"
	"merch_shop/pkg/middleware"
	tokenizermock "merch_shop/pkg/tokenizer/mocks"
	"merch_shop/pkg/xerrors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	service := New(dbmock.NewDB(t), slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

	t.Run("default invalid parms validation", func(t *testing.T) {
		testCases := []struct {
			name        string
			username    string
			password    string
			expectedErr xerrors.Xerror
		}{
			{name: "password lenth < min", username: strings.Repeat("1", minUsernameLenth),
				password: strings.Repeat("1", minPasswordLenth-1), expectedErr: xerrors.New(errPasswordInvalid, http.StatusBadRequest)},
			{name: "password lenth > max", username: strings.Repeat("1", minUsernameLenth),
				password: strings.Repeat("1", maxPasswordLenth+1), expectedErr: xerrors.New(errPasswordInvalid, http.StatusBadRequest)},
			{name: "username lenth < min", username: strings.Repeat("1", minUsernameLenth-1),
				password: strings.Repeat("1", minPasswordLenth), expectedErr: xerrors.New(errUsernameInvalid, http.StatusBadRequest)},
			{name: "username lenth > max", username: strings.Repeat("1", maxUsernameLenth+1),
				password: strings.Repeat("1", minPasswordLenth), expectedErr: xerrors.New(errUsernameInvalid, http.StatusBadRequest)},
		}

		for _, tc := range testCases {
			_, err := service.AuthentificateUser(context.Background(), tc.username, tc.password)
			assert.Equal(t, tc.expectedErr, err)
		}
	})

	username := strings.Repeat("1", minUsernameLenth)
	password := strings.Repeat("1", minPasswordLenth)
	userID := 1
	expToken := "test"

	t.Run("get user db error", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(nil, "", errors.New("some error"))

		_, err := service.AuthentificateUser(context.Background(), username, password)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("password encryption error", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(nil, "", db.ErrNoUser)
		cryptor.On("EncryptKeyword", mock.Anything).Return("", errors.New("some error"))

		_, err := service.AuthentificateUser(context.Background(), username, password)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("create user error", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(nil, "", db.ErrNoUser)
		datadase.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error"))
		cryptor.On("EncryptKeyword", mock.Anything).Return("", nil)

		_, err := service.AuthentificateUser(context.Background(), username, password)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})
	t.Run("password validation error", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(nil, "", nil)
		cryptor.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(errors.New("some error"))

		_, err := service.AuthentificateUser(context.Background(), username, password)
		require.Equal(t, xerrors.New(errPasswordMismatch, http.StatusUnauthorized), err)
	})

	t.Run("token generation error", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(&userID, "", nil)
		cryptor.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(nil)
		tokenizer.On("GenerateToken", mock.Anything).Return(nil, errors.New("some error"))

		_, err := service.AuthentificateUser(context.Background(), username, password)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("positive result with get user", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(&userID, "", nil)
		cryptor.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(nil)
		tokenizer.On("GenerateToken", mock.Anything).Return(&expToken, nil)

		token, err := service.AuthentificateUser(context.Background(), username, password)
		require.NoError(t, err)
		require.Equal(t, expToken, token)
	})

	t.Run("positive result with create user", func(t *testing.T) {
		datadase := dbmock.NewDB(t)
		tokenizer := tokenizermock.NewTokenizer(t)
		cryptor := cryptormock.NewCryptor(t)

		service := New(datadase, slog.Default(), cryptor, tokenizer)

		datadase.On("GetUser", mock.Anything, mock.Anything).Return(nil, "", db.ErrNoUser)
		datadase.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&userID, nil)
		cryptor.On("EncryptKeyword", mock.Anything).Return("", nil)
		tokenizer.On("GenerateToken", mock.Anything).Return(&expToken, nil)

		token, err := service.AuthentificateUser(context.Background(), username, password)
		require.NoError(t, err)
		require.Equal(t, expToken, token)
	})
}

func TestGetInfo(t *testing.T) {
	service := New(dbmock.NewDB(t), slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

	ctxEmpty := context.Background()
	ctxWithUserID := context.WithValue(ctxEmpty, middleware.UserIDKey, 1)
	balance := 1000

	t.Run("userID missing error", func(t *testing.T) {
		_, err := service.GetInfo(ctxEmpty)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("get info db error", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		database.On("GetUserInfoByUserID", mock.Anything, mock.Anything).Return(nil, nil, nil, errors.New("some error"))

		_, err := service.GetInfo(ctxWithUserID)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("positive result with empty inventory", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		expectedInfo := models.Info{
			Balance:         balance,
			Inventory:       []models.Item{},
			TransferHistory: models.CoinTransferHistory{},
		}

		database.On("GetUserInfoByUserID", mock.Anything, mock.Anything).Return(&balance, []byte(emptyJSONB), &models.CoinTransferHistory{}, nil)

		info, err := service.GetInfo(ctxWithUserID)
		require.NoError(t, err)
		require.Equal(t, expectedInfo, *info)
	})

	t.Run("positive result with non empty inventory", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		expectedInfo := models.Info{
			Balance: balance,
			Inventory: []models.Item{
				{
					Type:     "testing",
					Quantity: "100",
				},
				{
					Type:     "wewewe",
					Quantity: "1",
				},
			},
			TransferHistory: models.CoinTransferHistory{},
		}

		database.On("GetUserInfoByUserID", mock.Anything, mock.Anything).Return(
			&balance, []byte(`{"wewewe":"1"}, {"testing":"100"}`), &models.CoinTransferHistory{}, nil)

		info, err := service.GetInfo(ctxWithUserID)
		require.NoError(t, err)
		require.Equal(t, expectedInfo, *info)
	})
}

func TestBuyItem(t *testing.T) {
	service := New(dbmock.NewDB(t), slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

	ctxEmpty := context.Background()
	ctxWithUserID := context.WithValue(ctxEmpty, middleware.UserIDKey, 1)
	invalItemID := "invalid item id"
	validItemID := "1"

	t.Run("userID missing error", func(t *testing.T) {
		err := service.BuyItem(ctxEmpty, "")
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("itemID invalid error", func(t *testing.T) {
		err := service.BuyItem(ctxWithUserID, invalItemID)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("buy item db error", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		database.On("BuyItemByItemID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error"))

		err := service.BuyItem(ctxWithUserID, validItemID)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("buy item user side error", func(t *testing.T) {
		userErrors := []error{db.ErrNoItem, db.ErrNotEnoughCoins}

		for _, e := range userErrors {
			database := dbmock.NewDB(t)
			service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))
			database.On("BuyItemByItemID", mock.Anything, mock.Anything, mock.Anything).Return(e)

			err := service.BuyItem(ctxWithUserID, validItemID)
			require.Equal(t, xerrors.New(e, http.StatusBadRequest), err)
		}

	})

	t.Run("positive result", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		database.On("BuyItemByItemID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := service.BuyItem(ctxWithUserID, validItemID)
		require.NoError(t, err)
	})
}

func TestSendCoin(t *testing.T) {
	database := dbmock.NewDB(t)
	service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

	ctxEmpty := context.Background()
	ctxWithUserID := context.WithValue(ctxEmpty, middleware.UserIDKey, 1)
	invalidCoinAmountStr := minCoinsForTransfer - 1

	t.Run("coin amount error", func(t *testing.T) {
		err := service.SendCoin(ctxEmpty, "", invalidCoinAmountStr)
		require.Equal(t, xerrors.New(errCoinAmountInvalid, http.StatusBadRequest), err)
	})

	t.Run("userID missing error", func(t *testing.T) {
		err := service.SendCoin(ctxEmpty, "", minCoinsForTransfer)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("send coin db error", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		database.On("SendCoinByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error"))

		err := service.SendCoin(ctxWithUserID, "", minCoinsForTransfer)
		require.Equal(t, xerrors.New(errSmthWentWrong, http.StatusInternalServerError), err)
	})

	t.Run("send coin user side error", func(t *testing.T) {
		userErrors := []error{db.ErrNoUser, db.ErrNotEnoughCoins}

		for _, e := range userErrors {
			database := dbmock.NewDB(t)
			service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

			database.On("SendCoinByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e)

			err := service.SendCoin(ctxWithUserID, "", minCoinsForTransfer)
			require.Equal(t, xerrors.New(e, http.StatusBadRequest), err)
		}

	})

	t.Run("positive result", func(t *testing.T) {
		database := dbmock.NewDB(t)
		service := New(database, slog.Default(), cryptormock.NewCryptor(t), tokenizermock.NewTokenizer(t))

		database.On("SendCoinByUsername", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := service.SendCoin(ctxWithUserID, "", minCoinsForTransfer)
		require.NoError(t, err)
	})
}
