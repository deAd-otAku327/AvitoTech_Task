package db

import (
	"context"
	"merch_shop/internal/models"
)

func (s *storage) GetUserInfoByUserID(ctx context.Context, userID int) (*int, map[string]int, *models.CoinTransferHistory, error) {
	return nil, nil, nil, nil
}
