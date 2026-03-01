package client

import (
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
	Active   bool
}

func (t *Timer) Tick() (err error) {
	if !t.Active {
		return
	}
	for i, task := range t.Interval {
		if time.Since(task.LastTrigger) >= task.Interval {
			task.LastTrigger = time.Now()
			err = task.F()
			if err != nil {
				return
			}
		}
		t.Interval[i] = task
	}
	for i, task := range t.Once {
		if task.Trigger.After(time.Now()) {
			err = task.F()
			if err != nil {
				return
			}
			slices.Delete(t.Once, i, i)
		}
	}
	return
}

func (t *Timer) Every(interval time.Duration, f func() error) {
	t.Interval = append(t.Interval, IntervalTask{
		F:        f,
		Interval: interval,
	})
}

func (t *Timer) In(interval time.Duration, f func() error) {
	t.Once = append(t.Once, OnceTask{
		F:       f,
		Trigger: time.Now().Add(interval),
	})
}
