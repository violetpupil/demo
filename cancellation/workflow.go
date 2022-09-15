package cancellation

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

const Queue = "cancellation"
const WorkflowID = "cancellation-workflow"

func Workflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started")
	var a *Activities

	defer func() {
		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		err := workflow.ExecuteActivity(newCtx, a.Cleanup).Get(ctx, nil)
		if err != nil {
			logger.Error("cleanup activity error", "error", err)
		}
	}()

	var result string
	err := workflow.ExecuteActivity(ctx, a.Main).Get(ctx, &result)
	if err != nil {
		logger.Error("main activity error", "error", err)
		return err
	}
	logger.Info("main activity result", "result", result)
	logger.Info("workflow complete")
	return nil
}
