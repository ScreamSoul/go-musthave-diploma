package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/screamsoul/go-musthave-diploma/pkg/backoff"
)

func (r *PostgresRepository) GetWallet(ctx context.Context, userID uuid.UUID) (wallet *models.UserWallet, err error) {
	query := `SELECT balance, spent FROM loyalty_wallets WHERE user_id = $1`
	wallet = &models.UserWallet{}
	err = backoff.RetryWithBackoff(
		r.backoffInteraval,
		IsTemporaryConnectionError,
		r.logger,
		func() error {
			return r.db.GetContext(ctx, wallet, query, userID)
		},
	)

	if err != nil {
		return nil, err
	}
	return
}

func (r *PostgresRepository) WithdrawWallet(ctx context.Context, userID uuid.UUID, withdraw *models.Withdraw) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})

	if err != nil {
		return err
	}
	res, err := tx.ExecContext(
		ctx,
		`UPDATE loyalty_wallets SET balance = balance - $1, spent = spent + $2 WHERE user_id = $3 and balance > $4`,
		withdraw.Amount, withdraw.Amount, userID, withdraw.Amount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		tx.Rollback()
		return repositories.ErrLowBalance
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO loyalty_wallet_operations (order_number, user_id, amount) VALUES ($1, $2, $3)`,
		withdraw.Order, userID, withdraw.Amount,
	)

	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PostgresRepository) GetWithdrawals(ctx context.Context, userID uuid.UUID) (withdraws []models.Withdraw, err error) {
	query := `SELECT order_number, amount, processed_at FROM loyalty_wallet_operations WHERE user_id = $1 ORDER BY processed_at`
	err = backoff.RetryWithBackoff(
		r.backoffInteraval,
		IsTemporaryConnectionError,
		r.logger,
		func() error {
			return r.db.SelectContext(ctx, &withdraws, query, userID)
		},
	)

	if err != nil {
		return nil, err
	}
	return

}
