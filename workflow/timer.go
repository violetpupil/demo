package workflow

import (
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

// UpdatableTimer sleep until wake-up time.
// The wake-up time can be updated by channel.
// Inner timer will be replace with new when update.
type UpdatableTimer struct {
	Logger      log.Logger
	wakeUpTime  time.Time
	ctx         workflow.Context
	timerCtx    workflow.Context
	timerCancel workflow.CancelFunc
	timerFired  bool
}

func (u *UpdatableTimer) GetWakeUpTime() time.Time {
	return u.wakeUpTime
}

func (u *UpdatableTimer) SleepUntil(
	ctx workflow.Context,
	initialWakeUpTime time.Time,
	updateWakeUpTimeCh workflow.ReceiveChannel,
) {
	u.ctx = ctx
	u.wakeUpTime = initialWakeUpTime
	if u.Logger == nil {
		u.Logger = workflow.GetLogger(ctx)
	}

	// loop until timer fired or workflow is canceled
	// loop occur only wake-up time update
	for !u.timerFired && ctx.Err() == nil {
		u.timerCtx, u.timerCancel = workflow.WithCancel(ctx)
		duration := u.wakeUpTime.Sub(workflow.Now(u.timerCtx))
		timer := workflow.NewTimer(u.timerCtx, duration)
		u.Logger.Info("sleep...", "WakeUpTime", u.wakeUpTime)

		// block select one of processes
		workflow.NewSelector(u.timerCtx).
			AddFuture(timer, u.setTimerFired).
			AddReceive(updateWakeUpTimeCh, u.updateWakeUpTime).
			Select(u.timerCtx)
	}
}

// setTimerFired is called by selector, block for timer future ready:
// 1. arrived wake-up time
// 2. timer context canceled
// 3. workflow context canceled
func (u *UpdatableTimer) setTimerFired(f workflow.Future) {
	err := f.Get(u.timerCtx, nil)
	if err == nil {
		u.Logger.Info("timer fired")
		u.timerFired = true
	} else if u.ctx.Err() != nil {
		u.Logger.Info("workflow is canceled")
	}
}

// updateWakeUpTime is called when channel receive new wake-up time
func (u *UpdatableTimer) updateWakeUpTime(c workflow.ReceiveChannel, more bool) {
	// only cancel timer context
	u.timerCancel()
	c.Receive(u.timerCtx, &u.wakeUpTime)
	u.Logger.Info("wake up time updated")
}
