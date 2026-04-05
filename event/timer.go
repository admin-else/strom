package event

import (
	"context"
	"reflect"
	"slices"
	"time"
)

type IntervalTask struct {
	F           func() error
	LastTrigger time.Time
	Interval    time.Duration
}

type OnceTask struct {
	F       func() error
	Trigger time.Time
}

type Timer struct {
	Once     []OnceTask
	Interval []IntervalTask
}

func (t *Timer) Trigger(ch chan any) {
	triggerAt := time.Time{}
	found := false
	iOnce, iInterval := -1, -1

	for i, oTask := range t.Once {
		if oTask.Trigger.Before(triggerAt) || triggerAt.IsZero() {
			found = true
			triggerAt = oTask.Trigger
			iOnce = i
		}
	}

	for i, iTask := range t.Interval {
		if iTask.LastTrigger.Add(iTask.Interval).Before(triggerAt) || triggerAt.IsZero() {
			found = true
			triggerAt = iTask.LastTrigger.Add(iTask.Interval)
			iOnce = -1
			iInterval = i
		}
	}

	if !found {
		return
	}
	now := time.Now()
	sleepFor := triggerAt.Sub(now)
	//slog.Debug("sleep next", "dur", sleepFor)
	time.Sleep(sleepFor)
	if iOnce >= 0 {
		ch <- CallFunc{F: reflect.ValueOf(t.Once[iOnce].F), Args: []reflect.Value{}}
		t.Once = slices.Delete(t.Once, iOnce, iOnce+1)
		return
	}
	if iInterval >= 0 {
		t.Interval[iInterval].LastTrigger = triggerAt
		ch <- CallFunc{F: reflect.ValueOf(t.Interval[iInterval].F), Args: []reflect.Value{}}
		return
	}
	panic("unreachable")
}
func (t *Timer) Every(interval time.Duration, f func() error) {
	t.Interval = append(t.Interval, IntervalTask{
		F:           f,
		Interval:    interval,
		LastTrigger: time.Now(),
	})
}

func (t *Timer) In(interval time.Duration, f func() error) {
	t.Once = append(t.Once, OnceTask{
		F:       f,
		Trigger: time.Now().Add(interval),
	})
}

func (t *Timer) MakeEventSource(ctx context.Context) (chSending <-chan any) {
	ch := make(chan any)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				t.Trigger(ch)
			}
		}
	}()
	chSending = ch
	return
}

func (t *Timer) Start(c *Loop) {
	c.RegisterEventSource(t.MakeEventSource(c.Ctx))
}
