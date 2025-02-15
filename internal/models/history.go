package models

type CoinTransferHistory struct {
	Recieved []IngoingCoinTransfer  `json:"recieved"`
	Sent     []OutgoingCoinTransfer `json:"sent"`
}

type IngoingCoinTransfer struct {
	Username string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type OutgoingCoinTransfer struct {
	Username string `json:"toUser"`
	Amount   int    `json:"amount"`
}
