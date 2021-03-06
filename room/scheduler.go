package room

import (
	"time"
)

type Scheduler struct{}

// Creates a delayed task.
func (Scheduler) Delayed(delay time.Duration, fun func()) {
	go func() {
		ticker := time.NewTicker(delay)
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

// Creates a repeating task.
func (Scheduler) Repeating(period time.Duration, fun func(stop func())) func() {
	ch := make(chan bool)
	cancel := func() {
		ch <- true
	}
	go func() {
		ticker := time.NewTicker(period)
		for {
			select {
			case <-ticker.C:
				fun(cancel)
			case <-ch:
				return
			}
		}
	}()
	return cancel
}

// Creates a delayed repeating task.
func (*Scheduler) DelayedRepeating(delay time.Duration, period time.Duration, fun func(stop func())) func() {
	ch := make(chan bool)
	cancel := func() {
		ch <- true
	}
	go func() {
		ticker := time.NewTicker(delay)
		<-ticker.C
		ticker.Stop()
		ticker = time.NewTicker(period)
		for {
			select {
			case <-ticker.C:
				fun(cancel)
			case <-ch:
				return
			}
		}
	}()
	return cancel
}
