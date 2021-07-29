package room

import (
	"time"
)

type Scheduler struct{}

func (Scheduler) Delayed(duration time.Duration, fun func()) {
	go func() {
		ticker := time.NewTicker(duration)
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

func (Scheduler) Repeating(duration time.Duration, fun func()) func() {
	cancel := make(chan bool)
	go func() {
		ticker := time.NewTicker(duration)
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
