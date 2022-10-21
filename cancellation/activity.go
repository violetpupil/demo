package cancellation

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type Activities struct{}

func (a *Activities) Main(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("main activity started")
	for {
		select {
		case <-time.After(4 * time.Second):
			logger.Info("heartbeat recording")
			activity.RecordHeartbeat(ctx)
		case <-ctx.Done():
			logger.Info("context is canceled", "error", ctx.Err())
			return "I am canceled by Done", nil
		}
	}
}

func (a *Activities) Cleanup(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("cleanup activity started")
	return nil
}
