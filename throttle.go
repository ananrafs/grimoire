package grimoire

import (
	"sync/atomic"
	"time"
)

func throttle(duration time.Duration, sign <-chan struct{}, act func()) {
	lock := func(state *int32) {
		for !atomic.CompareAndSwapInt32(state, 0, 1) {
			time.Sleep(1 * time.Nanosecond)
		}
	}

	unlock := func(state *int32) {
		atomic.StoreInt32(state, 0)
	}

	state := int32(0)
	var timer <-chan time.Time
	for {
		select {
		case _, ok := <-sign:
			if !ok {
				return
			}
			lock(&state)
			if timer == nil {
				timer = time.After(duration)
			}
			unlock(&state)
		case <-timer:
			act()
			timer = nil
		}
	}
}
