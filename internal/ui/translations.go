package ui

import (
	"embed"
	"log/slog"

	"fyne.io/fyne/v2/lang"
)

//go:embed translations
var translationsFS embed.FS

func init() {
	if err := lang.AddTranslationsFS(translationsFS, "translations"); err != nil {
		slog.Error("failed to load translations", "error", err)
	}
}
