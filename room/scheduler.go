package room

import "time"

type Scheduler struct{}

func (Scheduler) Delayed(duration time.Duration, fun func()) {
	go func() {
		ticker := time.NewTicker(duration)
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}
