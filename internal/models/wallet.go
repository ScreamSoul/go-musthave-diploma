package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type UserWallet struct {
	Balance    float64 `json:"current" db:"balance"`
	TotalSpent float64 `json:"withdrawn" db:"spent"`
}

type Withdraw struct {
	Order       int64     `json:"order" db:"order_number"`
	Amount      float64   `json:"sum" db:"amount"`
	ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at,omitempty"`
}

func (w *Withdraw) MarshalJSON() ([]byte, error) {
	type Alias Withdraw
	aux := &struct {
		*Alias
		Order string `json:"order" db:"order_number"`
	}{
		Alias: (*Alias)(w),
		Order: func() string {
			return strconv.Itoa(int(w.Order))
		}(),
	}

	return json.Marshal(aux)
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
	orderNumber, err := strconv.ParseInt(aux.Order, 10, 64)
	if err != nil {
		return err
	}
	w.Order = orderNumber
	if w.Amount < 0 {
		return fmt.Errorf("the amount of the charge cannot be less than zero")
	}

	return nil
}
