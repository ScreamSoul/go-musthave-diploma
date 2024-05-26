package handlers

import (
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/pkg/logging"
	"go.uber.org/zap"
)

type UserLoyaltyServer struct {
	store  repositories.UserLoyaltyRepository
	logger *zap.Logger
}

func NewUserLoyaltyServer(loyaltyRepo repositories.UserLoyaltyRepository) *UserLoyaltyServer {
	logger := logging.GetLogger()

	return &UserLoyaltyServer{store: loyaltyRepo, logger: logger}
}
