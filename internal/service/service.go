package service

import (
	"merch_shop/internal/db"
	"merch_shop/internal/models"
)

type MerchShopService interface {
	Authentificate() error
	GetInfo() (*models.Info, error)
	BuyItem() (*models.Item, error)
	SendCoin() error
}

type merchShopService struct {
	storage db.DB
}

func New(storage db.DB) MerchShopService {
	return &merchShopService{storage: storage}
}

func (s *merchShopService) Authentificate() error {
	return nil
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
