package models

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

func (s OrderStatus) IsFinal() bool {
	switch s {
	case StatusInvalid, StatusProcessed:
		return true
	default:
		return false
	}
}
