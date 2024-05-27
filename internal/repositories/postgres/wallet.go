package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/pkg/backoff"
)

func (r *PostgresRepository) GetWalletInfo(ctx context.Context, userID uuid.UUID) (wallet *models.UserWallet, err error) {
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
