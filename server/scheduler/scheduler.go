package scheduler

import (
	"context"
	"time"
)

func StartScheduler(interval time.Duration, action func(ctx context.Context) error) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			action(ctx)
		}()

	}
}
