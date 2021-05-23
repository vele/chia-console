package chia

import (
	ui "github.com/gizak/termui/v3"
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
	ChiaPlotsEligableChart *widgets.BarChart
	ChiaProcessingTimes    *widgets.BarChart
}

func NewView() *View {
	var view = View{}

	view.Header = widgets.NewParagraph()
	view.Header.Border = false
	view.Header.Text = " Chia-console - Chia realtime inspector"
	view.Header.SetRect(0, 0, 0, 0)

	view.InfoBar = widgets.NewParagraph()
	view.InfoBar.Border = false
	view.InfoBar.Text = ""
	view.InfoBar.SetRect(0, 2, 0, 0)

	view.ChiaPlotsEligableChart = widgets.NewBarChart()
	view.ChiaPlotsEligableChart.Border = true
	view.ChiaPlotsEligableChart.Labels = []string{"Chia Eligable Plots"}
	view.ChiaPlotsEligableChart.BorderStyle.Fg = ui.ColorBlack
	view.ChiaPlotsEligableChart.SetRect(0, 2, 50, 0)

	return &view
}
