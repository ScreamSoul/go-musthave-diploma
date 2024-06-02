package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/pkg/backoff"
)

func (r *PostgresRepository) CreateUser(ctx context.Context, creds *models.Creds) (uuid.UUID, error) {

	userID := uuid.New()

	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	createUser := func() error {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3)`,
			userID, creds.Login, creds.Password)
		if err != nil {
			_ = tx.Rollback()
		}
		return err
	}
	var pgErr *pgconn.PgError
	err = backoff.RetryWithBackoff(r.backoffInteraval, IsTemporaryConnectionError, r.logger, createUser)
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = repositories.ErrUserAlreadyExists
	}

	if err != nil {
		return uuid.Nil, err
	}

	createLoyaltyWallet := func() error {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO loyalty_wallets (user_id, balance, spent) VALUES ($1, 0, 0)`,
			userID)
		if err != nil {
			_ = tx.Rollback()
		}
		return err
	}

	err = backoff.RetryWithBackoff(r.backoffInteraval, IsTemporaryConnectionError, r.logger, createLoyaltyWallet)
	if err != nil {
		return uuid.Nil, err
	}

	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func (r *PostgresRepository) CheckUserPassword(ctx context.Context, creds *models.Creds) (userID uuid.UUID, err error) {
	query := "SELECT id FROM users WHERE login = $1 and password_hash = $2"

	exec := func() error {
		return r.db.GetContext(ctx, &userID, query, creds.Login, creds.Password)
	}

	err = backoff.RetryWithBackoff(r.backoffInteraval, IsTemporaryConnectionError, r.logger, exec)
	if err == sql.ErrNoRows {
		err = repositories.ErrInvalidCredentials
	}

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *PostgresRepository) CreateOrder(ctx context.Context, orderNumber int, userID uuid.UUID) error {
	query := "SELECT COUNT(*) FROM orders WHERE number = $1 and user_id = $2"
	var count int
	checkExists := func() error {
		return r.db.GetContext(ctx, &count, query, orderNumber, userID)
	}

	err := backoff.RetryWithBackoff(r.backoffInteraval, IsTemporaryConnectionError, r.logger, checkExists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		return repositories.ErrOrderAlreadyUpload
	}

	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO orders (number, user_id)
		VALUES ($1, $2)
	`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	exec := func() error {
		_, err = stmt.ExecContext(ctx, orderNumber, userID)
		return err
	}

	err = backoff.RetryWithBackoff(r.backoffInteraval, IsTemporaryConnectionError, r.logger, exec)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = repositories.ErrOrderConflict
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) ListOrders(ctx context.Context, userID uuid.UUID) (orders []models.Order, err error) {
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1`
	exec := func() error {
		return r.db.SelectContext(ctx, &orders, query, userID)
	}

	err = backoff.RetryWithBackoff(r.backoffInteraval, IsTemporaryConnectionError, r.logger, exec)
	if err != nil {
		return nil, err
	}

	return
}
