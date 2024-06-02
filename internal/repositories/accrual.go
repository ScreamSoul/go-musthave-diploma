package repositories

import (
	"context"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
)

//go:generate minimock -i github.com/screamsoul/go-musthave-diploma/internal/repositories.AccrualRepository -o ../mocks/accural_repo_mock.go -g
type AccrualRepository interface {
	GetAccural(ctx context.Context, orderNumber int) (*models.Accural, error)
}
