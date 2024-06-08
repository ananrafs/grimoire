package grimoire

import (
	"testing"
	"time"
)

func Test_throttle(t *testing.T) {
	type args struct {
		duration time.Duration
		sign     <-chan struct{}
		act      func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			throttle(tt.args.duration, tt.args.sign, tt.args.act)
		})
	}
}
