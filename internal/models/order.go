package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Order struct {
	OrderNumber int           `json:"number" db:"number"`
	Status      OrderStatus   `json:"status" db:"status"`
	Accrual     sql.NullInt64 `json:"accrual,omitempty" db:"accrual,omitempty"`
	UploadedAt  time.Time     `json:"uploaded_at" db:"uploaded_at"`
}

func (u *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	aux := &struct {
		*Alias
		Accrual *int64 `json:"accrual,omitempty"`
	}{
		Alias: (*Alias)(u),
		Accrual: func() *int64 {
			var accural = &u.Accrual.Int64
			if u.Accrual.Valid {
				return accural
			}
			return nil
		}(),
	}

	return json.Marshal(aux)
}
