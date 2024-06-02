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
	logger           *zap.Logger
	backoffInteraval []time.Duration
}

func NewPostgresRepository(dataSourceName string, backoffInteraval []time.Duration) *PostgresRepository {
	db := sqlx.MustOpen("pgx", dataSourceName)

	return &PostgresRepository{db, logging.GetLogger(), backoffInteraval}
}

func (r *PostgresRepository) Ping(ctx context.Context) bool {
	err := r.db.PingContext(ctx)
	if err != nil {
		r.logger.Error("db connect error", zap.Error(err))
	}
	return err == nil
}

func (r *PostgresRepository) Close() {
	err := r.db.Close()
	if err != nil {
		r.logger.Error("db close connection error", zap.Error(err))
	}
}
