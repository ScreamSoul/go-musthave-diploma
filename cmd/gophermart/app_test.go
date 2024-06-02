package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Initializes Postgres repository with correct configuration

func TestAccural(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := &Config{
		Postgres: Postgres{
			DatabaseDSN:      "postgres://user:password@localhost/dbname",
			BackoffIntervals: []time.Duration{time.Second, 2 * time.Second},
		},
		JWT: JWT{
			Secret:          "secret",
			ExpiredDuration: time.Hour,
		},
		ActualSystemAddress: "http://localhost:8080",
	}
	logger := zap.NewExample()
	orderChain := make(chan int)

	assert.NotPanics(t, func() {
		go accural(ctx, cfg, logger, orderChain)
	})

	time.Sleep(time.Second)
}

func TestApp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := &Config{
		Postgres: Postgres{
			DatabaseDSN:      "postgres://user:password@localhost/dbname",
			BackoffIntervals: []time.Duration{time.Second, 2 * time.Second},
		},
		ActualSystemAddress: "http://localhost:8080",
	}
	logger := zap.NewExample()
	orderChain := make(chan int)

	assert.Panics(t, func() {
		app(ctx, cfg, logger, orderChain)
	})
}
