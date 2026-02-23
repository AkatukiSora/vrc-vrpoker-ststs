package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/application"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/applog"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/persistence"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/ui"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/watcher"
)

var (
	version   = "dev"
	commit    = "local"
	buildDate = "unknown"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	debug := *debugFlag || os.Getenv("VRC_VRPOKER_DEBUG") == "1"
	applog.Init(debug)

	slog.Info("starting",
		"version", version,
		"commit", commit,
		"buildDate", buildDate,
		"debug", debug,
	)

	dbPath := resolveDBPath()
	slog.Info("database", "path", dbPath)

	repo, err := persistence.NewSQLiteRepository(dbPath)
	if err != nil {
		slog.Warn("sqlite init failed, falling back to memory", "error", err)
		repo = nil
	}

	meta := ui.AppMetadata{
		Version:       version,
		Commit:        commit,
		BuildDate:     buildDate,
		RepositoryURL: "https://github.com/AkatukiSora/vrc-vrpoker-stats",
	}
	if repo != nil {
		ui.Run(application.NewService(repo, watcher.DetectAllLogFiles), meta, dbPath)
		return
	}

	ui.Run(application.NewService(persistence.NewMemoryRepository(), watcher.DetectAllLogFiles), meta, "")
}

// resolveDBPath returns the OS-appropriate path for the SQLite database:
//
//	Linux:   $XDG_DATA_HOME/vrc-vrpoker-ststs/vrpoker-stats.db
//	         (defaults to ~/.local/share/vrc-vrpoker-ststs/)
//	Windows: %LOCALAPPDATA%\vrc-vrpoker-ststs\vrpoker-stats.db
//	macOS:   ~/Library/Application Support/vrc-vrpoker-ststs/vrpoker-stats.db
//
// Falls back to ~/.vrc-vrpoker-ststs/ if the primary location is unavailable,
// then to the current directory as a last resort.
func resolveDBPath() string {
	const appName = "vrc-vrpoker-ststs"
	const dbFile = "vrpoker-stats.db"

	baseDir := userDataDir()
	if baseDir == "" {
		if home, err := os.UserHomeDir(); err == nil {
			dir := filepath.Join(home, "."+appName)
			if err := os.MkdirAll(dir, 0o755); err == nil {
				return filepath.Join(dir, dbFile)
			}
		}
		return filepath.Join(".", dbFile)
	}

	dir := filepath.Join(baseDir, appName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		slog.Warn("failed to create data directory", "dir", dir, "error", err)
		return filepath.Join(".", dbFile)
	}
	return filepath.Join(dir, dbFile)
}

// userDataDir returns the OS-specific base directory for persistent user data.
//
//	Linux:   $XDG_DATA_HOME, or ~/.local/share if unset
//	Windows: %LOCALAPPDATA%
//	macOS:   ~/Library/Application Support
func userDataDir() string {
	switch runtime.GOOS {
	case "windows":
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return dir
		}
		// Fallback: derive from USERPROFILE
		if profile := os.Getenv("USERPROFILE"); profile != "" {
			return filepath.Join(profile, "AppData", "Local")
		}
	case "darwin":
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "Library", "Application Support")
		}
	default: // Linux and other Unix-like
		if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
			return dir
		}
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".local", "share")
		}
	}
	return ""
}
