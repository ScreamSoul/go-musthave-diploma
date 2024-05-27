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
	ListOrders(ctx context.Context, userID uuid.UUID) ([]models.Order, error)

	UpdateOrderAccural(ctx context.Context, orderAccural *models.Accural) error

	GetWallet(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error)
	WithdrawWallet(ctx context.Context, userID uuid.UUID, withdraw *models.Withdraw) error
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]models.Withdraw, error)
}
