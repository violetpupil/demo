package workflow

import (
	"context"

	"go.temporal.io/sdk/activity"
)

type Activities struct{}

func (a *Activities) Main(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("activity done")
	return nil
}
