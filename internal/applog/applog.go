// Package applog initialises the global slog logger for the application.
// Call Init once at startup; all other packages use log/slog directly.
package applog

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var debugMode bool

// Init sets up the global slog logger.
// It writes structured text logs to both stdout and a temporary log file.
// If debug is true, the minimum log level is Debug; otherwise Info.
func Init(debug bool) {
	debugMode = debug

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	writers := []io.Writer{os.Stdout}
	if f, err := os.OpenFile(tempLogPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
		writers = append(writers, f)
	}

	h := slog.NewTextHandler(io.MultiWriter(writers...), &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(h))
}

// IsDebug reports whether debug mode is active.
func IsDebug() bool {
	return debugMode
}

func tempLogPath() string {
	return filepath.Join(os.TempDir(), "vrc-vrpoker-stats.log")
}
