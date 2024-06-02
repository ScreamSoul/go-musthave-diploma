package models

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

type Order struct {
	OrderNumber int64           `json:"number" db:"number"`
	Status      OrderStatus     `json:"status" db:"status"`
	Accrual     sql.NullFloat64 `json:"accrual,omitempty" db:"accrual,omitempty"`
	UploadedAt  time.Time       `json:"uploaded_at" db:"uploaded_at"`
}

func (u *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	aux := &struct {
		*Alias
		OrderNumber string   `json:"number" db:"number"`
		Accrual     *float64 `json:"accrual,omitempty"`
	}{
		Alias: (*Alias)(u),
		Accrual: func() *float64 {
			var accural = &u.Accrual.Float64
			if u.Accrual.Valid {
				return accural
			}
			return nil
		}(),
		OrderNumber: func() string {
			return strconv.Itoa(int(u.OrderNumber))
		}(),
	}

	return json.Marshal(aux)
}
