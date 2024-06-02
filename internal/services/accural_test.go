package services

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/screamsoul/go-musthave-diploma/internal/mocks"
	"github.com/screamsoul/go-musthave-diploma/internal/models"
	"github.com/screamsoul/go-musthave-diploma/internal/repositories"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AccuralUpdaterSuite struct {
	suite.Suite
	accuralRepo *mocks.AccrualRepositoryMock
	loyaltyRepo *mocks.UserLoyaltyRepositoryMock
	logger      *zap.Logger
}

func TestMemStorageSuite(t *testing.T) {
	suite.Run(t, new(AccuralUpdaterSuite))
}

func (s *AccuralUpdaterSuite) SetupTest() {
	mc := minimock.NewController(s.T())

	accuralRepo := mocks.NewAccrualRepositoryMock(mc)
	loyaltyRepo := mocks.NewUserLoyaltyRepositoryMock(mc)
	s.accuralRepo = accuralRepo
	s.loyaltyRepo = loyaltyRepo
	s.logger = zap.NewNop()
}

func (s *AccuralUpdaterSuite) TearDownTest() {
	s.loyaltyRepo.MinimockFinish()
	s.accuralRepo.MinimockFinish()
}

func (s *AccuralUpdaterSuite) TestStart() {

	// Настройка ожидаемого поведения моков
	orderNumber := 315
	status := models.StatusProcessed
	accuralValue := 500.0

	accural := &models.Accural{
		Order:   orderNumber,
		Status:  status,
		Accrual: &accuralValue,
	}

	var testTable = []struct {
		name            string
		mockAccuralRepo func(context.Context)
		mockLoyaltyRepo func(context.Context)
	}{
		{
			name: "default",
			mockAccuralRepo: func(ctx context.Context) {
				s.accuralRepo.GetAccuralMock.Expect(ctx, orderNumber).Return(accural, nil)
			},
			mockLoyaltyRepo: func(ctx context.Context) {
				s.loyaltyRepo.UpdateOrderAccuralMock.Expect(ctx, accural).Return(nil)
			},
		},
		{
			name: "invalid accural",
			mockAccuralRepo: func(ctx context.Context) {
				s.accuralRepo.GetAccuralMock.Expect(ctx, orderNumber).Return(nil, repositories.ErrInvalidAccural)
			},
			mockLoyaltyRepo: func(ctx context.Context) {
				s.loyaltyRepo.UpdateOrderAccuralMock.Expect(ctx, &models.Accural{Status: models.StatusInvalid, Accrual: nil, Order: orderNumber}).Return(nil)
			},
		},
	}

	for _, v := range testTable {
		s.Suite.Run(v.name, func() {
			var wg sync.WaitGroup
			wg.Add(1)

			orderChain := make(chan int)
			accuralUpdater := NewOrderAccuralUpdater(
				s.loyaltyRepo,
				s.accuralRepo,
				s.logger,
				orderChain,
			)

			ctx, cancel := context.WithCancel(context.Background())

			v.mockAccuralRepo(ctx)
			v.mockLoyaltyRepo(ctx)
			go accuralUpdater.Start(ctx)

			go func() {
				defer wg.Done()
				accuralUpdater.Start(ctx)
				close(orderChain)
			}()

			// Отправляем тестовый заказ в канал
			orderChain <- orderNumber

			time.Sleep(time.Second)
			cancel()
			wg.Wait()
		})

	}
}

// TODO: add processing recursion call when order have status pricessing or new
