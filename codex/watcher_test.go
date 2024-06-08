package codex

import "testing"

func Test_watchChanges(t *testing.T) {
	type args struct {
		loc        string
		onModified func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watchChanges(tt.args.loc, tt.args.onModified)
		})
	}
}
