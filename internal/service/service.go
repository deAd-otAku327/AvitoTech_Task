package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"merch_shop/internal/db"
	"merch_shop/internal/models"
	"merch_shop/pkg/cryptor"
	"merch_shop/pkg/middleware"
	"merch_shop/pkg/tokenizer"
	"merch_shop/pkg/xerrors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	minCoinsForTransfer = 1

	minPasswordLength = 6
	maxPasswordLength = 15

	minUsernameLength = 1
	maxUsernameLength = 10

	emptyJSONB = "{}"
)

var (
	errSmthWentWrong    = errors.New("something went wrong")
	errPasswordMismatch = errors.New("incorrect password")

	errCoinAmountInvalid = fmt.Errorf("coin amount is invalid: min %d", minCoinsForTransfer)
	errPasswordInvalid   = fmt.Errorf("password is invalid: length min %d max %d", minPasswordLength, maxPasswordLength)
	errUsernameInvalid   = fmt.Errorf("username is invalid: length min %d max %d", minUsernameLength, maxUsernameLength)
)

var jsonbFormatRunesToCut = map[rune]struct{}{
	'{': {},
	'"': {},
	' ': {},
	'}': {},
}

type MerchShopService interface {
	AuthentificateUser(ctx context.Context, username, password string) (string, xerrors.Xerror)
	GetInfo(ctx context.Context) (*models.Info, xerrors.Xerror)
	BuyItem(ctx context.Context, itemID string) xerrors.Xerror
	SendCoin(ctx context.Context, destUsername string, amount int) xerrors.Xerror
}

type merchShopService struct {
	storage   db.DB
	logger    *slog.Logger
	cryptor   cryptor.Cryptor
	tokenizer tokenizer.Tokenizer
}

func New(storage db.DB, log *slog.Logger, cr cryptor.Cryptor, t tokenizer.Tokenizer) MerchShopService {
	return &merchShopService{
		storage:   storage,
		logger:    log,
		cryptor:   cr,
		tokenizer: t,
	}
}

func (s *merchShopService) AuthentificateUser(ctx context.Context, username, password string) (string, xerrors.Xerror) {
	if len(password) > maxPasswordLength || len(password) < minPasswordLength {
		return "", xerrors.New(errPasswordInvalid, http.StatusBadRequest)
	}
	if len(username) > maxUsernameLength || len(username) < minUsernameLength {
		return "", xerrors.New(errUsernameInvalid, http.StatusBadRequest)
	}

	userID, dbPassword, err := s.storage.GetUser(ctx, username)
	if err != nil {
		if err == db.ErrNoUser {
			encryptedPass, err := s.cryptor.EncryptKeyword(password)
			if err != nil {
				s.logger.Error("encrypt password: " + err.Error())
				return "", xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
			}

			userID, err = s.storage.CreateUser(ctx, username, encryptedPass)
			if err != nil {
				s.logger.Error("create user: " + err.Error())
				return "", xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
			}
		} else {
			s.logger.Error("get user: " + err.Error())
			return "", xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
		}

	} else {
		if err := s.cryptor.CompareHashAndPassword(dbPassword, password); err != nil {
			return "", xerrors.New(errPasswordMismatch, http.StatusUnauthorized)
		}
	}

	token, err := s.tokenizer.GenerateToken(strconv.Itoa(*userID))
	if err != nil {
		s.logger.Error("generate token: " + err.Error())
		return "", xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	return *token, nil
}

func (s *merchShopService) GetInfo(ctx context.Context) (*models.Info, xerrors.Xerror) {
	userID, ok := ctx.Value(middleware.UserIDKey).(int)
	if !ok {
		return nil, xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	balance, inventory, history, err := s.storage.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("get user info: " + err.Error())
		return nil, xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	proccessedInventory := make([]models.Item, 0)
	if string(inventory) != emptyJSONB {

		cleanInventory := strings.Builder{}
		for _, r := range string(inventory) {
			if _, ok := jsonbFormatRunesToCut[r]; !ok {
				cleanInventory.WriteRune(r)
			}
		}

		splittedInventory := strings.Split(cleanInventory.String(), ",")
		for _, val := range splittedInventory {
			typeToAmountPair := strings.Split(val, ":")
			proccessedInventory = append(proccessedInventory, models.Item{
				Type:     typeToAmountPair[0],
				Quantity: typeToAmountPair[1],
			})
		}

		sort.Slice(proccessedInventory, func(i, j int) bool {
			return proccessedInventory[i].Type < proccessedInventory[j].Type
		})

	}

	info := &models.Info{
		Balance:         *balance,
		Inventory:       proccessedInventory,
		TransferHistory: *history,
	}

	return info, nil
}

func (s *merchShopService) BuyItem(ctx context.Context, itemIDStr string) xerrors.Xerror {
	userID, ok := ctx.Value(middleware.UserIDKey).(int)
	if !ok {
		return xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		s.logger.Error("itemID string to int: " + err.Error())
		return xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	err = s.storage.BuyItemByItemID(ctx, userID, itemID)
	if err != nil {
		if err == db.ErrNoItem || err == db.ErrNotEnoughCoins {
			return xerrors.New(err, http.StatusBadRequest)
		}
		s.logger.Error("buy item: " + err.Error())
		return xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}
	return nil
}

func (s *merchShopService) SendCoin(ctx context.Context, destUsername string, amount int) xerrors.Xerror {
	if amount < minCoinsForTransfer {
		return xerrors.New(errCoinAmountInvalid, http.StatusBadRequest)
	}
	userID, ok := ctx.Value(middleware.UserIDKey).(int)
	if !ok {
		return xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	err := s.storage.SendCoinByUsername(ctx, userID, destUsername, amount)
	if err != nil {
		if err == db.ErrNoUser || err == db.ErrNotEnoughCoins {
			return xerrors.New(err, http.StatusBadRequest)
		}
		s.logger.Error("send coin: " + err.Error())
		return xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	return nil
}
