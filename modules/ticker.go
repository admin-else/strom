package modules

import "time"

type IntervalTask struct {
	F           func() error
	LastTrigger time.Time
	Interval    time.Duration
}

type OnceTask struct {
	F       func() error
	Trigger time.Time
}

type Ticker struct {
	Once     []OnceTask
	Interval []IntervalTask
	Active   bool
}

func (t *Ticker) Tick() (err error) {
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
	var triggerLater []OnceTask
	for _, task := range t.Once {
		if task.Trigger.After(time.Now()) {
			err = task.F()
			if err != nil {
				return
			}
		} else {
			triggerLater = append(triggerLater, task)
		}
	}
	t.Once = triggerLater
	return
}
