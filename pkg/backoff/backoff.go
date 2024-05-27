package backoff

import (
	"time"

	"go.uber.org/zap"
)

func RetryWithBackoff(
	backoffIntervals []time.Duration,
	shouldRetry func(error) bool,
	logger *zap.Logger,
	fn func() error,
) error {
	var err error

	if len(backoffIntervals) == 0 {
		return fn()
	}

	for i := 0; i < len(backoffIntervals); i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if shouldRetry(err) {
			if i < len(backoffIntervals) {
				logger.Warn("retry with err", zap.Error(err))
				time.Sleep(backoffIntervals[i])
			}
		} else {
			logger.Warn("failed retries with err", zap.Error(err))
			return err
		}

	}
	return err
}
