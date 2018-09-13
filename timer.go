package util

import (
	"context"
	"log"
	"runtime/debug"
	"time"
)

// CallbackFunc callback function
type CallbackFunc func()

// DailyTask daily execution
func DailyTask(ctx context.Context, hour, minute, second int, loc *time.Location,
	callback CallbackFunc) {
	go func(ctx context.Context, hour, minute, second int, loc *time.Location,
		callback CallbackFunc) {
		// prevent exec callback panic
		defer func() {
			if err := recover(); err != nil {
				log.Printf("exec callback panic, %v\n", err)
				log.Printf("stack:%s", debug.Stack())
				// rerun DailyTask
				go DailyTask(ctx, hour, minute, second, loc, callback)
			}
		}()
		for {

			now := time.Now().In(loc)
			timing := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, loc)
			// obtain the difference between timing and now
			diff := timing.Sub(now)
			if diff < 0 {
				// add 1 day
				diff = timing.AddDate(0, 0, 1).Sub(now)
			}
			select {
			case <-time.After(diff):
				go callback()
			case <-ctx.Done():
				break
			}
		}
	}(ctx, hour, minute, second, loc, callback)
}

// DailyTaskCST DailyTask for CST
func DailyTaskCST(ctx context.Context, hour, minute, second int, callback CallbackFunc) {
	DailyTask(ctx, hour, minute, second, CST(), callback)
}

// DailyTaskUTC DailyTask for UTC
func DailyTaskUTC(ctx context.Context, hour, minute, second int, callback CallbackFunc) {
	DailyTask(ctx, hour, minute, second, time.UTC, callback)
}
