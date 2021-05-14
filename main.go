package main

import (
	"flag"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var (
	certFile = flag.String("cert", "/root/.chia/mainnet/config/ssl/full_node/private_full_node.crt", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "/root/.chia/mainnet/config/ssl/full_node/private_full_node.key", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "/root/.chia/mainnet/config/ssl/ca/chia_ca.crt", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()
	os.Setenv("CHIA_CA_CRT", *caFile)
	os.Setenv("CHIA_FULL_NODE_CRT", *certFile)
	os.Setenv("CHIA_FULL_NODE_KEY", *keyFile)
	os.Setenv("CHIA_CONSOLE_CONFIGURED", "1")
	os.Setenv("CHIA_SERVER_URL", "https://127.0.0.1:8555")

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
