package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type UserWallet struct {
	Balance    int `json:"current" db:"balance"`
	TotalSpent int `json:"withdrawn" db:"spent"`
}

type Withdraw struct {
	Order  int `json:"order"`
	Amount int `json:"sum"`
}

func (w *Withdraw) UnmarshalJSON(data []byte) error {
	type Alias Withdraw
	aux := &struct {
		Order string `json:"order"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	orderNumber, err := strconv.Atoi(aux.Order)
	if err != nil {
		return err
	}
	w.Order = orderNumber
	if w.Amount < 0 {
		return fmt.Errorf("the amount of the charge cannot be less than zero")
	}

	return nil
}
