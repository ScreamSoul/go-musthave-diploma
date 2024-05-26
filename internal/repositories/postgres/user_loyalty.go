package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/pkg/backoff"
	"github.com/screamsoul/go-musthave-diploma/pkg/logging"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db               *sqlx.DB
	logging          *zap.Logger
	backoffInteraval []time.Duration
}

func NewPostgresStorage(dataSourceName string, backoffInteraval []time.Duration) *PostgresStorage {
	db := sqlx.MustOpen("pgx", dataSourceName)

	return &PostgresStorage{db, logging.GetLogger(), backoffInteraval}
}

func (storage *PostgresStorage) Ping(ctx context.Context) bool {
	err := storage.db.PingContext(ctx)
	if err != nil {
		storage.logging.Error("db connect error", zap.Error(err))
	}
	return err == nil
}

func (storage *PostgresStorage) Close() {
	err := storage.db.Close()
	if err != nil {
		storage.logging.Error("db close connection error", zap.Error(err))
	}
}

func (storage *PostgresStorage) CreateUser(ctx context.Context, creds *models.Creds) (uuid.UUID, error) {
	userID := uuid.New()
	stmt, err := storage.db.PrepareContext(ctx, `
		INSERT INTO users (id, login, password_hash)
		VALUES ($1, $2, $3)
	`)

	if err != nil {
		return uuid.Nil, err
	}
	defer stmt.Close()

	exec := func() error {
		_, err = stmt.ExecContext(ctx, userID, creds.Login, creds.Password)
		return err
	}

	err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = repositories.ErrUserAlreadyExists
	} else if err != nil {
		err = fmt.Errorf("failed retries db request, %w", err)
	}

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func (storage *PostgresStorage) CheckUserPassword(ctx context.Context, creds *models.Creds) (userID uuid.UUID, err error) {
	query := "SELECT id FROM users WHERE login = $1 and password_hash = $2"

	exec := func() error {
		return storage.db.GetContext(ctx, &userID, query, creds.Login, creds.Password)
	}

	err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
	if err == sql.ErrNoRows {
		err = repositories.ErrInvalidCredentials
	} else if err != nil {
		err = fmt.Errorf("failed retries db request, %w", err)
	}

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (storage *PostgresStorage) CreateOrder(ctx context.Context, orderNumber int, userId uuid.UUID) error {
	query := "SELECT COUNT(*) FROM orders WHERE number = $1 and user_id = $2"
	var count int
	check_exists := func() error {
		return storage.db.GetContext(ctx, &count, query, orderNumber, userId)
	}

	err := backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, check_exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed retries db request, %w", err)
	}
	if count > 0 {
		return repositories.ErrOrderAlreadyUpload
	}

	stmt, err := storage.db.PrepareContext(ctx, `
		INSERT INTO orders (number, user_id)
		VALUES ($1, $2)
	`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	exec := func() error {
		_, err = stmt.ExecContext(ctx, orderNumber, userId)
		return err
	}

	err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = repositories.ErrOrderConflict
	} else if err != nil {
		err = fmt.Errorf("failed retries db request, %w", err)
	}

	if err != nil {
		return err
	}

	return nil
}

func (storage *PostgresStorage) ListOrders(ctx context.Context, userID uuid.UUID) (orders []models.Order, err error) {
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1`
	exec := func() error {
		return storage.db.SelectContext(ctx, &orders, query, userID)
	}

	err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
	if err != nil {
		err = fmt.Errorf("failed retries db request, %w", err)
	}

	return
}
