package models

type History struct {
	Recieved []CoinTransfer
	Sent     []CoinTransfer
}

type CoinTransfer struct {
	Username string
	Amount   int
}
