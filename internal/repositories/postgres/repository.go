package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/screamsoul/go-musthave-diploma/pkg/logging"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	db               *sqlx.DB
	logging          *zap.Logger
	backoffInteraval []time.Duration
}

func NewPostgresRepository(dataSourceName string, backoffInteraval []time.Duration) *PostgresRepository {
	db := sqlx.MustOpen("pgx", dataSourceName)

	return &PostgresRepository{db, logging.GetLogger(), backoffInteraval}
}

func (storage *PostgresRepository) Ping(ctx context.Context) bool {
	err := storage.db.PingContext(ctx)
	if err != nil {
		storage.logging.Error("db connect error", zap.Error(err))
	}
	return err == nil
}

func (storage *PostgresRepository) Close() {
	err := storage.db.Close()
	if err != nil {
		storage.logging.Error("db close connection error", zap.Error(err))
	}
}
