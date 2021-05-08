package main

import (
	"flag"
	"fmt"
	"log"

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
	result, err := chiaClient.GetChiaBlockchainState("127.0.0.1:8555")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

}
