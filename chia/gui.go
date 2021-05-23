package chia

import (
	"github.com/gizak/termui/v3/widgets"
)

type UIEvent int

const (
	KeyArrowUp UIEvent = 1 << iota
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyCtrlC
	KeyCtrlD
	KeyQ
	Resize
	KeyI
)

type View struct {
	Header                 *widgets.Paragraph
	InfoBar                *widgets.Paragraph
	ChiaPlotsEligableChart *widgets.SparklineGroup
	ChiaProcessingTimes    *widgets.BarChart
}

func NewView() *View {
	var view = View{}

	return &view
}
