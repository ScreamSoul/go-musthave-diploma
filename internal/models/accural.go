package models

import (
	"encoding/json"
	"strconv"
)

type Accural struct {
	Order   int         `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual int         `json:"accrual"`
}

func (u *Accural) UnmarshalJSON(data []byte) error {
	type Alias Accural
	aux := &struct {
		Order  string `json:"order"`
		Status string `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Status == "REGISTERED" {
		u.Status = StatusNew
	} else {
		u.Status = OrderStatus(aux.Status)
	}

	orderNumber, err := strconv.Atoi(aux.Order)

	if err != nil {
		return err
	}

	u.Order = orderNumber

	return nil
}
