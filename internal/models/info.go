package models

type Info struct {
	Balance         int
	Inventory       []Item
	TransferHistory History
}
