package codex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestWatchChanges(t *testing.T) {
	tests := []struct {
		name        string
		createFiles []string
		modifyFiles []string
		expectCalls int
	}{
		{
			name:        "single event",
			createFiles: []string{"file1.txt"},
			modifyFiles: []string{"file1.txt"},
			expectCalls: 1,
		},
		{
			name:        "multiple events",
			createFiles: []string{"file2.txt", "file3.txt"},
			modifyFiles: []string{"file2.txt", "file3.txt"},
			expectCalls: 2,
		},
		{
			name:        "no events",
			createFiles: []string{"file4.txt"},
			modifyFiles: []string{},
			expectCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := ioutil.TempDir("", "watcher_test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			for _, file := range tt.createFiles {
				path := filepath.Join(tempDir, file)
				if err := ioutil.WriteFile(path, []byte("initial content"), 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			var wg sync.WaitGroup
			wg.Add(tt.expectCalls)

			onModified := func() {
				wg.Done()
			}

			go watchChanges(tempDir, onModified)

			time.Sleep(500 * time.Millisecond) // Give watcher time to set up

			for _, file := range tt.modifyFiles {
				path := filepath.Join(tempDir, file)
				if err := ioutil.WriteFile(path, []byte("modified content"), 0644); err != nil {
					t.Fatalf("Failed to modify file %s: %v", file, err)
				}
			}

			// Wait for all expected calls to be made
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Fatalf("Test timed out, expected %d calls but not all were made", tt.expectCalls)
			}
		})
	}
}
