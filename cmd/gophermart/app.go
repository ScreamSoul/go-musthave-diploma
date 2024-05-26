package main

import (
	"context"
	"net/http"

	"github.com/screamsoul/go-musthave-diploma/internal/handlers"
	"github.com/screamsoul/go-musthave-diploma/internal/middlewares"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories/postgres"
	"github.com/screamsoul/go-musthave-diploma/internal/routers"
	"github.com/screamsoul/go-musthave-diploma/internal/services"
	"go.uber.org/zap"
)

func app(ctx context.Context, cfg *Config, logger *zap.Logger) error {

	logger.Info("init postgres")
	postgresStorage := postgres.NewPostgresStorage(cfg.DatabaseDSN, cfg.BackoffIntervals)
	defer postgresStorage.Close()

	logger.Info("start migrations")
	if err := postgresStorage.Bootstrap(ctx); err != nil {
		return err
	}

	logger.Info("init token service")
	services.Initialize(cfg.JWT.Secret, cfg.JWT.ExpiredDuration)

	logger.Info("init user loyalty server")
	var loyaltyServer = handlers.NewUserLoyaltyServer(
		postgresStorage,
	)

	logger.Info("init user loyalty router")
	var router = routers.NewUserLoyaltyRouter(
		loyaltyServer,
		middlewares.LoggingMiddleware,
		middlewares.GzipDecompressMiddleware,
		middlewares.GzipCompressMiddleware,
	)

	logger.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))
	return http.ListenAndServe(cfg.ListenAddress, router)
}
