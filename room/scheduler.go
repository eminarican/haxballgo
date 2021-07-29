package room

import (
	"time"
)

type Scheduler struct{}

func (Scheduler) Delayed(delay time.Duration, fun func()) {
	go func() {
		ticker := time.NewTicker(delay)
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

func (Scheduler) Repeating(period time.Duration, fun func()) func() {
	cancel := make(chan bool)
	go func() {
		ticker := time.NewTicker(period)
		for {
			select {
			case <-ticker.C:
				fun()
		    case <-cancel:
				return
			}
		}
	}()
	return func() {
		cancel <- true
	}
}
