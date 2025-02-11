package models

type User struct {
	Username string
	Password string

	Balance         int
	Inventory       []Item
	TransferHistory History
}
