package repositories

import "context"

type AccrualPointsRepository interface {
	Ping(ctx context.Context) bool
	GetPoints(ctx context.Context, orderNumber int) (int, error)
}
