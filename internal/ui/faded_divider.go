package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const (
	dividerFadeStartFromEdge = 0.30
	dividerFadeEndFromEdge   = 0.10
)

var dividerBaseColor = color.NRGBA{R: 0xAC, G: 0xAF, B: 0xB5, A: 0xFF}

func newSectionDivider() fyne.CanvasObject {
	r := canvas.NewRasterWithPixels(func(x, _x, w, _h int) color.Color {
		if w <= 1 {
			return dividerBaseColor
		}
		t := float32(x) / float32(w-1)
		edge := t
		if 1-t < edge {
			edge = 1 - t
		}
		alphaFactor := float32(1)
		if edge <= dividerFadeEndFromEdge {
			alphaFactor = 0
		} else if edge < dividerFadeStartFromEdge {
			alphaFactor = (edge - dividerFadeEndFromEdge) / (dividerFadeStartFromEdge - dividerFadeEndFromEdge)
		}
		c := dividerBaseColor
		c.A = uint8(float32(c.A) * alphaFactor)
		return c
	})
	r.SetMinSize(fyne.NewSize(0, 1))
	return r
}

func newSectionDividerVertical() fyne.CanvasObject {
	r := canvas.NewRasterWithPixels(func(_x, y, _w, h int) color.Color {
		if h <= 1 {
			return dividerBaseColor
		}
		t := float32(y) / float32(h-1)
		edge := t
		if 1-t < edge {
			edge = 1 - t
		}
		alphaFactor := float32(1)
		if edge <= dividerFadeEndFromEdge {
			alphaFactor = 0
		} else if edge < dividerFadeStartFromEdge {
			alphaFactor = (edge - dividerFadeEndFromEdge) / (dividerFadeStartFromEdge - dividerFadeEndFromEdge)
		}
		c := dividerBaseColor
		c.A = uint8(float32(c.A) * alphaFactor)
		return c
	})
	r.SetMinSize(fyne.NewSize(1, theme.TextSize()*1.8))
	return r
}
