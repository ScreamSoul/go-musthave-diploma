package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
)

func (r *PostgresRepository) UpdateOrderAccural(ctx context.Context, orderAccural *models.Accural) error {

	tx, err := r.db.BeginTxx(ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		},
	)

	if err != nil {
		return err
	}

	// update accural and status in orders by number and get user_id
	var res sql.Result
	res, err = tx.ExecContext(
		ctx, `UPDATE orders SET accrual = $1, status = $2 WHERE number = $3`,
		orderAccural.Accrual, orderAccural.Status, orderAccural.Order,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Print(res.RowsAffected())

	if orderAccural.Accrual == 0 || orderAccural.Status != models.StatusProcessed {
		return tx.Commit()
	}

	// if accural > 0 and status processed
	// get user_id by order_number
	var userID uuid.UUID
	err = tx.QueryRow(`SELECT user_id FROM orders WHERE number = $1`, orderAccural.Order).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update loyalty_wallets bu user_id
	_, err = tx.ExecContext(
		ctx, `UPDATE loyalty_wallets SET balance = balance + $1 WHERE user_id = $2`,
		orderAccural.Accrual, userID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
