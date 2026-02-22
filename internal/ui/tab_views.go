package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/stats"
)

type overviewTabView struct {
	root       *fyne.Container
	win        fyne.Window
	visibility *MetricVisibilityState
}

func newOverviewTabView(win fyne.Window, visibility *MetricVisibilityState) *overviewTabView {
	return &overviewTabView{root: container.NewMax(), win: win, visibility: visibility}
}

func (v *overviewTabView) Update(s *stats.Stats) {
	v.root.Objects = []fyne.CanvasObject{NewOverviewTab(s, v.visibility, v.win)}
	v.root.Refresh()
}

func (v *overviewTabView) CanvasObject() fyne.CanvasObject {
	return v.root
}

type positionStatsTabView struct {
	root       *fyne.Container
	visibility *MetricVisibilityState
}

func newPositionStatsTabView(visibility *MetricVisibilityState) *positionStatsTabView {
	return &positionStatsTabView{root: container.NewMax(), visibility: visibility}
}

func (v *positionStatsTabView) Update(s *stats.Stats) {
	v.root.Objects = []fyne.CanvasObject{NewPositionStatsTab(s, v.visibility)}
	v.root.Refresh()
}

func (v *positionStatsTabView) CanvasObject() fyne.CanvasObject {
	return v.root
}
