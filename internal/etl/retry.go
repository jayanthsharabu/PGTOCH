package etl

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      bool
}

func Retry(ctx context.Context, config RetryConfig, operation func() error) error {
	var attempt int
	for {
		err := operation()
		if err == nil {
			return nil
		}
		attempt++
		if attempt >= config.MaxAttempts {
			return errors.New("max Attempts Reached")
		}

		backoff := min(config.BaseDelay*time.Duration(math.Pow(2, float64(attempt))), config.MaxDelay)
		if config.Jitter {
			jitterRange := float64(backoff) * 0.25
			jitter := time.Duration(rand.Float64()*jitterRange*2 - jitterRange)
			backoff += jitter

			if backoff < 0 {
				backoff = config.BaseDelay
			}
		}

		zap.L().Warn("Operation failed, retrying...",
			zap.Int("attempt", attempt),
			zap.Error(err),
			zap.Duration("backoff", backoff),
		)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			continue
		}
	}
}
