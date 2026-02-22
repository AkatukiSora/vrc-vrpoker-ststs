package main

import (
	"path/filepath"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/application"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/persistence"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/ui"
	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/watcher"
)

func main() {
	dbPath := filepath.Join(".", "vrpoker-stats.db")
	repo, err := persistence.NewSQLiteRepository(dbPath)
	if err != nil {
		repo = nil
	}

	if repo != nil {
		ui.Run(application.NewService(repo, watcher.DetectAllLogFiles))
		return
	}

	ui.Run(application.NewService(persistence.NewMemoryRepository(), watcher.DetectAllLogFiles))
}
