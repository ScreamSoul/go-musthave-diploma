package repositories

import (
	"context"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
)

type AccrualRepository interface {
	GetAccural(ctx context.Context, orderNumber int) (*models.Accural, error)
}
