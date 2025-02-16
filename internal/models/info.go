package models

type Info struct {
	Balance         int                 `json:"coins"`
	Inventory       []Item              `json:"inventory"`
	TransferHistory CoinTransferHistory `json:"coinHistory"`
}
