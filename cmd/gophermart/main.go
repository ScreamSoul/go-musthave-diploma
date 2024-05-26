package main

import (
	"context"

	"github.com/screamsoul/go-musthave-diploma/pkg/logging"
)

func main() {
	// Init context
	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	// Init Config
	cfg, err := NewConfig()

	if err != nil {
		panic(err)
	}

	// Init logger
	if err := logging.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger := logging.GetLogger()

	if err := app(ctx, cfg, logger); err != nil {
		panic(err)
	}
}
