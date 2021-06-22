package main

import (
	"flag"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
	"github.com/vele/chia-console/chia"
)

var (
	certFile          = flag.String("cert", "/root/.flax/mainnet/config/ssl/full_node/private_full_node.crt", "A PEM eoncoded certificate file.")
	keyFile           = flag.String("key", "/root/.flax/mainnet/config/ssl/full_node/private_full_node.key", "A PEM encoded private key file.")
	walletCertFile    = flag.String("walletCrt", "/root/.flax/mainnet/config/ssl/wallet/private_wallet.crt", "A PEM encoded private key file.")
	walletKeyFile     = flag.String("walletKey", "/root/.flax/mainnet/config/ssl/wallet/private_wallet.key", "A PEM encoded private key file.")
	harvesterCertFile = flag.String("harvesterCrt", "/root/.flax/mainnet/config/ssl/harvester/private_harvester.crt", "A PEM encoded private key file.")
	harvesterKeyFile  = flag.String("harvesterKey", "/root/.flax/mainnet/config/ssl/harvester/private_harvester.key", "A PEM encoded private key file.")
	caFile            = flag.String("CA", "/root/.flax/mainnet/config/ssl/ca/chia_ca.crt", "A PEM eoncoded CA's certificate file.")
	logFile           = flag.String("LogFileDir", " /root/.flax/mainnet/log/debug.log", "The location of chia log files.")
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
	os.Setenv("CHIA_SERVER_URL", "https://127.0.0.1:6755")
	os.Setenv("CHIA_WALLET_URL", "https://127.0.0.1:6883")
	os.Setenv("CHIA_HARVESTER_URL", "https://127.0.0.1:6760")
	os.Setenv("CHIA_LOGFILE", *logFile)
	err := chia.CreateSqlSchema()
	if err != nil {
		log.Panicln(err)
	}
	err = chia.FetchChiaMap()
	if err != nil {
		log.Panicln(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(firstRowGraph)
	g.SupportOverlaps = false

	wg.Add(1)
	go drawEligablePlotsGraph(g)
	go drawProcessingTimesGraph(g)
	go drawFreeSpaceTable(g)
	go updateChiaPriceDB()
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
