package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/pkg/backoff"
)

func (storage *PostgresRepository) CreateUser(ctx context.Context, creds *models.Creds) (uuid.UUID, error) {

	userID := uuid.New()

	tx, err := storage.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	create_user := func() error {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3)`,
			userID, creds.Login, creds.Password)
		if err != nil {
			tx.Rollback()
		}
		return err
	}
	var pgErr *pgconn.PgError
	err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, create_user)
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = repositories.ErrUserAlreadyExists
	} else if err != nil {
		err = fmt.Errorf("failed retries db request, %w", err)
	}

	if err != nil {
		return uuid.Nil, err
	}

	create_loyalty_wallet := func() error {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO loyalty_wallets (user_id, balance, spent_points_total) VALUES ($1, 0, 0)`,
			userID)
		if err != nil {
			tx.Rollback()
		}
		return err
	}

	err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, create_loyalty_wallet)
	if err != nil {
		return uuid.Nil, err
	}

	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func (storage *PostgresRepository) CheckUserPassword(ctx context.Context, creds *models.Creds) (userID uuid.UUID, err error) {
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

func (storage *PostgresRepository) CreateOrder(ctx context.Context, orderNumber int, userId uuid.UUID) error {
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

func (storage *PostgresRepository) ListOrders(ctx context.Context, userID uuid.UUID) (orders []models.Order, err error) {
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
