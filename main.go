package main

import (
	"flag"
	"log"
	"os"

	"github.com/jroimartin/gocui"
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
	logFile           = flag.String("LogFileDir", " /root/.chia/mainnet/log/debug.log", "The location of chia log files.")
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
	err := chia.CreateSqlSchema()
	if err != nil {
		log.Panicln(err)
	}
	err = chia.FetchChiaMap()
	if err != nil {
		log.Panicln(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(mainLayout)
	wg.Add(1)
	go drawEligablePlotsGraph(g)
	go drawProcessingTimesGraph(g)
	go drawFreeSpaceTable(g)
	go updateChiaPriceDB(g)
	go updateChiaPriceGUI(g)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
