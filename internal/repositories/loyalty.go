package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
)

//go:generate minimock -i github.com/screamsoul/go-musthave-diploma/internal/repositories.UserLoyaltyRepository -o ../mocks/loyalty_repo_mock.go -g

type UserLoyaltyRepository interface {
	Ping(ctx context.Context) bool
	CreateUser(ctx context.Context, creds *models.Creds) (uuid.UUID, error)
	CheckUserPassword(ctx context.Context, creds *models.Creds) (uuid.UUID, error)

	CreateOrder(ctx context.Context, orderNumber int, userID uuid.UUID) error
	ListOrders(ctx context.Context, userID uuid.UUID) ([]models.Order, error)

	UpdateOrderAccural(ctx context.Context, orderAccural *models.Accural) error

	GetWallet(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error)
	WithdrawWallet(ctx context.Context, userID uuid.UUID, withdraw *models.Withdraw) error
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]models.Withdraw, error)
}
