package main

import (
	"flag"
	"log"
	"os"

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

	var SparkLineData []float64
	fetchLogs := chia.ParseLogs()
	for item := range fetchLogs {
		SparkLineData = append(SparkLineData, float64(fetchLogs[item].Plots))
	}
	ChiaPlotsSparkline := w.NewSparkline()
	ChiaPlotsSparkline.Data = SparkLineData
	ChiaPlotsEligableChart := w.NewSparklineGroup(ChiaPlotsSparkline)
	ChiaPlotsEligableChart.Title = "Eligable Plot Counts"
	ChiaPlotsEligableChart.BorderStyle.Fg = ui.ColorBlue

	//ChiaPlotsEligableChart.SetRect(0, 0, 100, 40)
	grid.Set(
		ui.NewRow(0.1/1,
			ui.NewCol(1.0/1, header),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, ChiaPlotsEligableChart),
		),
	)
	ui.Render(grid)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
