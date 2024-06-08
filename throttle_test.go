package grimoire

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_throttle(t *testing.T) {
	type tests struct {
		name              string
		duration          time.Duration
		signals           int
		signalInterval    time.Duration
		expectedCallCount int32
		waitTime          time.Duration
	}

	testCases := []tests{
		{
			name:              "Multi call due to throttling",
			duration:          50 * time.Millisecond,
			signals:           10,
			signalInterval:    10 * time.Millisecond,
			expectedCallCount: 2,
			waitTime:          200 * time.Millisecond,
		},
		{
			name:              "Multiple calls without throttling",
			duration:          50 * time.Millisecond,
			signals:           5,
			signalInterval:    60 * time.Millisecond,
			expectedCallCount: 5,
			waitTime:          300 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var callCount int32
			sign := make(chan struct{})
			wg := sync.WaitGroup{}

			act := func() {
				atomic.AddInt32(&callCount, 1)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				throttle(tc.duration, sign, act)
			}()

			// Send signals based on the test case
			for i := 0; i < tc.signals; i++ {
				sign <- struct{}{}
				time.Sleep(tc.signalInterval)
			}

			// Wait enough time to ensure act() could have been called
			time.Sleep(tc.waitTime)

			// Close the sign channel to stop the throttle goroutine
			close(sign)

			// Wait for the throttle goroutine to finish
			wg.Wait()

			// Verify the callCount
			if callCount != tc.expectedCallCount {
				t.Errorf("Expected act() to be called %d times, but it was called %d times", tc.expectedCallCount, callCount)
			}
		})
	}
}
