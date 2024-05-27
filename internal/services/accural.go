package services

import (
	"context"
	"errors"
	"time"

	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"go.uber.org/zap"
)

type OrderAccuralUpdater struct {
	loyaltyRepo repositories.UserLoyaltyRepository
	accuralRepo repositories.AccrualRepository
	logger      *zap.Logger
	orderChain  chan int
}

func NewOrderAccuralUpdater(
	loyaltyRepo repositories.UserLoyaltyRepository,
	accuralRepo repositories.AccrualRepository,
	logger *zap.Logger,
	orderChain chan int,
) *OrderAccuralUpdater {
	return &OrderAccuralUpdater{
		loyaltyRepo: loyaltyRepo,
		accuralRepo: accuralRepo,
		logger:      logger,
		orderChain:  orderChain,
	}
}

func (s *OrderAccuralUpdater) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("close accural updater worker")
			return
		case orderNumber := <-s.orderChain:
			s.logger.Info("start check accural on order", zap.Int("order", orderNumber))
			accutalOrder, err := s.accuralRepo.GetAccural(ctx, orderNumber)
			if errors.Is(err, repositories.ErrInvalidAccural) {
				accutalOrder = &models.Accural{Status: models.StatusInvalid, Accrual: nil, Order: orderNumber}
			} else if err != nil {
				s.logger.Error("Fail GetAccural", zap.Error(err))
				break
			}

			err = s.loyaltyRepo.UpdateOrderAccural(
				ctx,
				accutalOrder,
			)

			if err != nil {
				s.logger.Error("Fail UpdateOrderAccural", zap.Error(err))
				break
			}

			switch accutalOrder.Status {
			case models.StatusInvalid, models.StatusProcessed:
				s.logger.Info("calculate accural finish", zap.Int("order", orderNumber))
				break
			default:
				s.logger.Info("wait for the calculations accural to be completed", zap.Int("order", orderNumber))

				go func() {
					time.Sleep(5 * time.Second)
					s.orderChain <- orderNumber
				}()
			}
		}
	}
}
