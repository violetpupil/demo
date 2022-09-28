package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

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
