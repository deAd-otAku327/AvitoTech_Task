package service

import (
	"context"
	"errors"
	"merch_shop/internal/db"
	"merch_shop/internal/models"
	"merch_shop/pkg/cryptor"
	"merch_shop/pkg/tokenizer"
)

type MerchShopService interface {
	AuthentificateUser(ctx context.Context, username, password string) (string, error)
	GetInfo() (*models.Info, error)
	BuyItem() (*models.Item, error)
	SendCoin() error
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

	token, err := s.tokenizer.GenerateToken(*userID)
	if err != nil {
		return "", err
	}

	return *token, nil
}

func (s *merchShopService) GetInfo() (*models.Info, error) {
	return nil, nil
}

func (s *merchShopService) BuyItem() (*models.Item, error) {
	return nil, nil
}

func (s *merchShopService) SendCoin() error {
	return nil
}
