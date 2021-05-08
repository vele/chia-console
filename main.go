package main

import (
	"flag"
	"fmt"

	"github.com/vele/chia-console/chia"
)

var (
	certFile = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()
	chiaClient := chia.NewClient(*certFile, *keyFile, *caFile)
	fmt.Println(chiaClient)

}
