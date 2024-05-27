package models

type UserWallet struct {
	Balance    int `json:"current" db:"balance"`
	TotalSpent int `json:"withdrawn" db:"spent"`
}
