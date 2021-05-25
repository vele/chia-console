package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
	"github.com/vele/chia-console/chia"
)

var (
	certFile          = flag.String("cert", "/root/.chia/mainnet/config/ssl/full_node/private_full_node.crt", "A PEM eoncoded certificate file.")
	keyFile           = flag.String("key", "/root/.chia/mainnet/config/ssl/full_node/private_full_node.key", "A PEM encoded private key file.")
	walletCertFile    = flag.String("walletCrt", "/root/.chia/mainnet/config/ssl/wallet/private_wallet.crt", "A PEM encoded private key file.")
	walletKeyFile     = flag.String("walletKey", "/root/.chia/mainnet/config/ssl/wallet/private_wallet.key", "A PEM encoded private key file.")
	harvesterCertFile = flag.String("harvesterCrt", "/root/.chia/mainnet/config/ssl/harvester/private_harvester.crt", "A PEM encoded private key file.")
	harvesterKeyFile  = flag.String("harvesterKey", "/root/.chia/mainnet/config/ssl/harvester/private_harvester.key", "A PEM encoded private key file.")
	caFile            = flag.String("CA", "/root/.chia/mainnet/config/ssl/ca/chia_ca.crt", "A PEM eoncoded CA's certificate file.")
	logFile           = flag.String("LogFile", " /root/.chia/mainnet/log/debug.log", "The location of chia log files.")
)

func main() {
	flag.Parse()
	os.Setenv("CHIA_CA_CRT", *caFile)
	os.Setenv("CHIA_FULL_NODE_CRT", *certFile)
	os.Setenv("CHIA_FULL_NODE_KEY", *keyFile)
	os.Setenv("CHIA_WALLET_CRT", *walletCertFile)
	os.Setenv("CHIA_WALLET_KEY", *walletKeyFile)
	os.Setenv("CHIA_HARVESTER_CRT", *harvesterCertFile)
	os.Setenv("CHIA_HARVESTER_KEY", *harvesterKeyFile)
	os.Setenv("CHIA_SERVER_URL", "https://127.0.0.1:8555")
	os.Setenv("CHIA_WALLET_URL", "https://127.0.0.1:9256")
	os.Setenv("CHIA_HARVESTER_URL", "https://127.0.0.1:8560")
	os.Setenv("CHIA_LOGFILE", *logFile)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	header := w.NewParagraph()
	header.Border = false
	header.Text = " Chia-console - Chia realtime inspector"
	header.SetRect(0, 0, 0, 0)

	var PlotCounters []float64
	//now := time.Now()
	fetchLogs := chia.ParseLogs()
	for item := range fetchLogs {
		PlotCounters = append(PlotCounters, float64(fetchLogs[item].Plots))
	}
	//log.Println(SparkLineData)
	ChiaPlotsSparkline := w.NewPlot()
	ChiaPlotsSparkline.Data = make([][]float64, 1)
	ChiaPlotsSparkline.Data[0] = PlotCounters
	ChiaPlotsSparkline.AxesColor = ui.ColorWhite
	ChiaPlotsSparkline.LineColors[0] = ui.ColorGreen
	ChiaPlotsSparkline.Title = "Eligable Plot Counts"
	ChiaPlotsSparkline.BorderStyle.Fg = ui.ColorBlue
	ChiaPlotsSparkline.TitleStyle.Fg = ui.ColorYellow
	ChiaPlotsSparkline.TitleStyle.Bg = ui.ColorBlack

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(0.5/1, ChiaPlotsSparkline),
		),
	)
	draw := func() {
		getPlotCounters := PopulateLogData()
		ChiaPlotsSparkline.Data[0] = getPlotCounters
		ChiaPlotsSparkline.Title = fmt.Sprintf("Eligable Plot Counts %s ", time.Now().String())
		ui.Render(grid)
	}
	draw()
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second * 5).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			draw()
		}
	}
}
func PopulateLogData() []float64 {
	var PlotCounters []float64
	fetchLogs := chia.ParseLogs()
	for item := range fetchLogs {
		PlotCounters = append(PlotCounters, float64(fetchLogs[item].Plots))
	}
	return PlotCounters
}
