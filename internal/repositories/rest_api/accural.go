package restapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
)

func (r *AccuralAPIRepository) GetAccural(ctx context.Context, orderNumber int) (accural *models.Accural, err error) {
	urlPath := fmt.Sprintf("/api/orders/%d", orderNumber)

	resp, err := r.request(ctx, "GET", urlPath, models.Accural{})

	if resp.StatusCode() == http.StatusNoContent {
		return nil, repositories.ErrInvalidAccural
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, repositories.ErrManyRequests
	}

	if resp.StatusCode() > 200 || err != nil {
		return nil, fmt.Errorf("bad request; %w", err)
	}

	accural, ok := resp.Result().(*models.Accural)
	if !ok {
		return nil, repositories.ErrInvalidAccural
	}

	return
}
