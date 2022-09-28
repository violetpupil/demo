package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Context = workflow.Context

const Queue = "demo-queue"
const WorkflowID = "demo-workflow"
const QueryType = "GetWakeUpTime"
const SignalType = "UpdateWakeUpTime"

func Workflow(ctx workflow.Context, initialWakeUpTime time.Time) error {
	logger := workflow.GetLogger(ctx)
	timer := UpdatableTimer{}
	err := workflow.SetQueryHandler(ctx, QueryType, func() (time.Time, error) {
		return timer.GetWakeUpTime(), nil
	})
	if err != nil {
		return err
	}
	timer.SleepUntil(ctx, initialWakeUpTime, workflow.GetSignalChannel(ctx, SignalType))

	logger.Info("wake up, execute activity")
	var a *Activities
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	future := workflow.ExecuteActivity(ctx, a.Main)
	err = future.Get(ctx, nil)
	return err
}

func DelayWorkflow(ctx workflow.Context, wakeUpTime time.Time) error {
	logger := workflow.GetLogger(ctx)

	timer := workflow.NewTimer(ctx, wakeUpTime.Sub(workflow.Now(ctx)))
	logger.Info("sleep...", "WakeUpTime", wakeUpTime)
	workflow.NewSelector(ctx).AddFuture(timer, func(f workflow.Future) {}).Select(ctx)

	logger.Info("wake up, execute activity")
	var a *Activities
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	future := workflow.ExecuteActivity(ctx, a.Main)
	err := future.Get(ctx, nil)
	return err
}
