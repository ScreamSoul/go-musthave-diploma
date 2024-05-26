package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
)

type UserLoyaltyRepository interface {
	Ping(ctx context.Context) bool
	CreateUser(ctx context.Context, creds *models.Creds) (uuid.UUID, error)
	CheckUserPassword(ctx context.Context, creds *models.Creds) (uuid.UUID, error)

	CreateOrder(ctx context.Context, orderNumber int, userId uuid.UUID) error
	ListOrders(ctx context.Context, userID uuid.UUID) (orders []models.Order, err error)
}
