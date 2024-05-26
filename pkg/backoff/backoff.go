package backoff

import "time"

func RetryWithBackoff(
	backoffIntervals []time.Duration,
	shouldRetry func(error) bool,
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
				time.Sleep(backoffIntervals[i])
			}
		} else {
			return err
		}

	}
	return err
}
