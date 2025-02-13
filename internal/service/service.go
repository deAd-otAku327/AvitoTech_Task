package service

import (
	"context"
	"errors"
	"merch_shop/internal/db"
	"merch_shop/internal/models"
	"merch_shop/pkg/cryptor"
	"merch_shop/pkg/middleware"
	"merch_shop/pkg/tokenizer"
	"strconv"
)

type MerchShopService interface {
	AuthentificateUser(ctx context.Context, username, password string) (string, error)
	GetInfo() (*models.Info, error)
	BuyItem(ctx context.Context, itemID string) error
	SendCoin(ctx context.Context, destUsername string, amount int) error
}

var ErrPasswordMismatch = errors.New("invalid password")

type merchShopService struct {
	tokenizer *tokenizer.Tokenizer
	storage   db.DB
}

func New(storage db.DB, t *tokenizer.Tokenizer) MerchShopService {
	return &merchShopService{
		storage:   storage,
		tokenizer: t,
	}
}

func (s *merchShopService) AuthentificateUser(ctx context.Context, username, password string) (string, error) {
	encryptedPass, err := cryptor.EncryptKeyword(password)
	if err != nil {
		return "", err
	}

	userID, dbPassword, err := s.storage.CreateOrGetUser(ctx, username, encryptedPass)
	if err != nil {
		return "", err
	}

	if err := cryptor.CompareHashAndPassword(dbPassword, password); err != nil {
		return "", ErrPasswordMismatch
	}

	token, err := s.tokenizer.GenerateToken(strconv.Itoa(*userID))
	if err != nil {
		return "", err
	}

	return *token, nil
}

func (s *merchShopService) GetInfo() (*models.Info, error) {
	return nil, nil
}

func (s *merchShopService) BuyItem(ctx context.Context, itemIDStr string) error {
	userID := ctx.Value(middleware.UserIDKey).(int)
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return err
	}

	err = s.storage.BuyItemByItemID(ctx, userID, itemID)
	if err != nil {
		return err
	}
	return nil
}

func (s *merchShopService) SendCoin(ctx context.Context, destUsername string, amount int) error {
	userID := ctx.Value(middleware.UserIDKey).(int)

	err := s.storage.SendCoinByUsername(ctx, userID, destUsername, amount)
	if err != nil {
		return err
	}

	return nil
}
