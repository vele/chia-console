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

	view.Header = widgets.NewParagraph()
	view.Header.Border = false
	view.Header.Text = " Chia-console - Chia realtime inspector"
	view.Header.SetRect(0, 0, 0, 0)

	var SparkLineData []float64
	fetchLogs := ParseLogs()
	for item := range fetchLogs {
		SparkLineData = append(SparkLineData, float64(fetchLogs[item].Plots))
	}
	ChiaPlotsSparkline := widgets.NewSparkline()
	ChiaPlotsSparkline.Data = SparkLineData
	view.ChiaPlotsEligableChart = widgets.NewSparklineGroup(ChiaPlotsSparkline)
	view.ChiaPlotsEligableChart.Title = "Sparkline 0"
	view.ChiaPlotsEligableChart.SetRect(0, 0, 20, 10)

	return &view
}
