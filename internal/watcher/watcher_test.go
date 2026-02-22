package watcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogWatcherStopIsIdempotent(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "output_log_2026-02-21_00-00-00.txt")
	if err := os.WriteFile(logPath, []byte(""), 0o600); err != nil {
		t.Fatalf("write log: %v", err)
	}

	lw, err := NewLogWatcher(logPath)
	if err != nil {
		t.Fatalf("new watcher: %v", err)
	}

	lw.Stop()
	lw.Stop()
}
