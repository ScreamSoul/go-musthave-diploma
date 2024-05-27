package main

import (
	"time"

	"github.com/alexflint/go-arg"
)

type Postgres struct {
	DatabaseDSN      string          `arg:"-d,env:DATABASE_URI" default:"" help:"Строка подключения к базе Postgres"`
	BackoffIntervals []time.Duration `arg:"--b-intervals,env:BACKOFF_INTERVALS" help:"Интервалы повтора запроса (обязательно если (default=1s,3s,5s)"`
	BackoffRetries   bool            `arg:"--backoff,env:BACKOFF_RETRIES" default:"true" help:"Повтор запроса при разрыве соединения"`
}

type JWT struct {
	Secret          string        `arg:"-a,env:JWT_SECRET" default:"qwerty123" help:"JWT секретный ключ"`
	ExpiredDuration time.Duration `arg:"-a,env:JWT_EXPIRED" default:"1h" help:"Время жизни токена"`
}

type Config struct {
	Postgres
	JWT
	ListenAddress       string `arg:"-a,env:RUN_ADDRESS" default:"localhost:8080" help:"Адрес и порт сервера"`
	LogLevel            string `arg:"--ll,env:LOG_LEVEL" default:"INFO" help:"Уровень логирования"`
	ActualSystemAddress string `arg:"-r,env:ACCRUAL_SYSTEM_ADDRESS" default:"http://localhost:8001" help:"Адрес системы расчёта начислений"`
	AccuralCheckerLimit int    `arg:"-r,env:ACCRUAL_CHECKER_LIMIT" default:"5" help:"Количество одновременных обращений к системе расчёта начислений"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := arg.Parse(&cfg); err != nil {
		return nil, err
	}

	if cfg.Postgres.BackoffIntervals == nil && cfg.Postgres.BackoffRetries {
		cfg.Postgres.BackoffIntervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	} else if !cfg.Postgres.BackoffRetries {
		cfg.Postgres.BackoffIntervals = nil
	}

	return &cfg, nil
}
