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

type MerchShopService interface {
	AuthentificateUser(ctx context.Context, username, password string) (string, xerrors.Xerror)
	GetInfo(ctx context.Context) (*models.Info, xerrors.Xerror)
	BuyItem(ctx context.Context, itemID string) xerrors.Xerror
	SendCoin(ctx context.Context, destUsername string, amount int) xerrors.Xerror
}

const (
	minCoinsForTransfer = 1

	minPasswordLenth = 6
	maxPasswordLenth = 15

	minUsernameLenth = 1
	maxUsernameLenth = 10
)

var (
	errSmthWentWrong    = errors.New("something went wrong")
	errPasswordMismatch = errors.New("incorrect password")

	errCoinAmountInvalid = fmt.Errorf("coin amount is invalid: min %d", minCoinsForTransfer)
	errPasswordInvalid   = fmt.Errorf("password is invalid: lenth min %d max %d", minPasswordLenth, maxPasswordLenth)
	errUsernameInvalid   = fmt.Errorf("username is invalid: lenth min %d max %d", minUsernameLenth, maxUsernameLenth)
)

type merchShopService struct {
	storage   db.DB
	logger    *slog.Logger
	tokenizer tokenizer.Tokenizer
}

func New(storage db.DB, log *slog.Logger, t tokenizer.Tokenizer) MerchShopService {
	return &merchShopService{
		storage:   storage,
		logger:    log,
		tokenizer: t,
	}
}

func (s *merchShopService) AuthentificateUser(ctx context.Context, username, password string) (string, xerrors.Xerror) {
	if len(password) > maxPasswordLenth || len(password) < minPasswordLenth {
		return "", xerrors.New(errPasswordInvalid, http.StatusBadRequest)
	}
	if len(username) > maxUsernameLenth || len(username) < minUsernameLenth {
		return "", xerrors.New(errUsernameInvalid, http.StatusBadRequest)
	}

	userID, dbPassword, err := s.storage.GetUser(ctx, username)
	if err != nil {
		if err == db.ErrNoUser {

			encryptedPass, err := cryptor.EncryptKeyword(password)
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
		if err := cryptor.CompareHashAndPassword(dbPassword, password); err != nil {
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
	userID := ctx.Value(middleware.UserIDKey).(int)

	balance, inventory, history, err := s.storage.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("get user info: " + err.Error())
		return nil, xerrors.New(errSmthWentWrong, http.StatusInternalServerError)
	}

	proccessedInventory := make([]models.Item, 0)
	if string(inventory) != "{}" {

		cleanInventory := strings.Builder{}
		for _, r := range string(inventory) {
			if r != '{' && r != '}' && r != '"' && r != ' ' {
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
	userID := ctx.Value(middleware.UserIDKey).(int)
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
	userID := ctx.Value(middleware.UserIDKey).(int)

	err := s.storage.SendCoinByUsername(ctx, userID, destUsername, amount)
	if err != nil {
		if err == db.ErrNoUser || err == db.ErrNotEnoughCoins {
			return xerrors.New(err, http.StatusBadRequest)
		}
		s.logger.Error("send coin: " + err.Error())
		return xerrors.New(err, http.StatusInternalServerError)
	}

	return nil
}
