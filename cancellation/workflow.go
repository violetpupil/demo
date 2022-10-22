package cancellation

import (
	"errors"
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const Queue = "cancellation"
const WorkflowID = "cancellation-workflow"

func GetExecuteLogger(ctx workflow.Context, activity string) log.Logger {
	logger := workflow.GetLogger(ctx)
	return log.With(logger, "Activity", activity)
}

func Workflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started")
	ExecuteMain(ctx)
	if errors.Is(ctx.Err(), workflow.ErrCanceled) {
		logger.Info("workflow canceled")
	} else if ctx.Err() != nil {
		logger.Error("workflow error", "Error", ctx.Err())
	} else {
		logger.Info("workflow complete")
	}
	return nil
}

func ExecuteMain(ctx workflow.Context) {
	// prepare
	logger := GetExecuteLogger(ctx, "main")
	logger.Info("execute activity")
	var a *Activities

	// execute
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var result string
	err := workflow.ExecuteActivity(ctx, a.Main).Get(ctx, &result)

	// process get return
	actualError := errors.Unwrap(err)
	switch actualError.(type) {
	case *temporal.TimeoutError:
		logger.Error("activity error, timeout", "Error", err)
	case *temporal.CanceledError:
		logger.Error("activity error, canceled", "Error", err)
	case *temporal.ApplicationError:
		logger.Error("activity error, application", "Error", err)
	case *temporal.PanicError:
		logger.Error("activity error, panic", "Error", err)
	case nil:
		logger.Info("activity complete", "Result", result)
	default:
		logger.Error("activity error, unknown", "Error", err)
	}
	ExecuteCleanup(ctx)
}

func ExecuteCleanup(ctx workflow.Context) {
	// prepare
	logger := GetExecuteLogger(ctx, "cleanup")
	logger.Info("execute activity")
	var a *Activities

	// execute
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	ctx, _ = workflow.NewDisconnectedContext(ctx)
	err := workflow.ExecuteActivity(ctx, a.Cleanup).Get(ctx, nil)
	if err != nil {
		logger.Error("activity error", "Error", err)
	} else {
		logger.Info("activity complete")
	}
}
