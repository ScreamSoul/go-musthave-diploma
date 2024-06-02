package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/screamsoul/go-musthave-diploma/internal/handlers"
	"github.com/screamsoul/go-musthave-diploma/internal/middlewares"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories/postgres"
	restapi "github.com/screamsoul/go-musthave-diploma/internal/repositories/rest_api"
	"github.com/screamsoul/go-musthave-diploma/internal/routers"
	"github.com/screamsoul/go-musthave-diploma/internal/services"
	"go.uber.org/zap"
)

func app(ctx context.Context, cfg *Config, logger *zap.Logger, orderChain chan int) {
	logger.Info("init postgres")
	loyaltyRepository := postgres.NewPostgresRepository(cfg.DatabaseDSN, cfg.BackoffIntervals)
	defer loyaltyRepository.Close()

	logger.Info("start migrations")
	if err := loyaltyRepository.Bootstrap(ctx); err != nil {
		panic(err)
	}

	logger.Info("init token service")
	services.Initialize(cfg.JWT.Secret, cfg.JWT.ExpiredDuration)

	logger.Info("init user loyalty server")
	var loyaltyServer = handlers.NewUserLoyaltyServer(
		loyaltyRepository,
		orderChain,
	)

	logger.Info("init user loyalty router")
	var router = routers.NewUserLoyaltyRouter(
		loyaltyServer,
		middlewares.LoggingMiddleware,
		middlewares.GzipDecompressMiddleware,
		middlewares.GzipCompressMiddleware,
	)

	logger.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}
}

func accural(ctx context.Context, cfg *Config, logger *zap.Logger, orderChain chan int) {
	logger.Info("init postgres")
	loyaltyRepository := postgres.NewPostgresRepository(cfg.DatabaseDSN, cfg.BackoffIntervals)
	defer loyaltyRepository.Close()

	logger.Info("Init accural api repository")
	accuralRepository := restapi.NewAccuralAPIRepository(cfg.ActualSystemAddress)

	updater := services.NewOrderAccuralUpdater(loyaltyRepository, accuralRepository, logger, orderChain)
	updater.Start(ctx)
}

func start(ctx context.Context, cfg *Config, logger *zap.Logger) {
	orderChain := make(chan int)

	go app(ctx, cfg, logger, orderChain)

	for i := 0; i < cfg.AccuralCheckerLimit; i++ {
		go accural(ctx, cfg, logger, orderChain)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return
	case <-sigChan:
		return
	}
}
