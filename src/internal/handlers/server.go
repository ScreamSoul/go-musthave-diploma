package handlers

import (
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/pkg/logging"
	"go.uber.org/zap"
)

type UserLoyaltyServer struct {
	loyaltyRepo repositories.UserLoyaltyRepository
	logger      *zap.Logger
	orderChain  chan int
}

func NewUserLoyaltyServer(loyaltyRepo repositories.UserLoyaltyRepository, orderChain chan int) *UserLoyaltyServer {
	logger := logging.GetLogger()

	return &UserLoyaltyServer{loyaltyRepo: loyaltyRepo, orderChain: orderChain, logger: logger}
}
