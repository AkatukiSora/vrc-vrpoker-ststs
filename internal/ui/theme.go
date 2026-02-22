package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// defaultPokerTheme is a dark theme for the poker stats app
type defaultPokerTheme struct{}

var _ fyne.Theme = (*defaultPokerTheme)(nil)

func (t defaultPokerTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	base := theme.DarkTheme().Color(name, theme.VariantDark)
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x3F, G: 0x81, B: 0xC6, A: 0xFF}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x66, G: 0xA8, B: 0xE0, A: 0xAA}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x7F, G: 0x8D, B: 0x9B, A: 0x2A}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x49, G: 0x89, B: 0xCA, A: 0x44}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x20, G: 0x26, B: 0x2D, A: 0xFF}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0x23, G: 0x2A, B: 0x32, A: 0xFF}
	default:
		return base
	}
}

func (t defaultPokerTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func (t defaultPokerTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (t defaultPokerTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}
