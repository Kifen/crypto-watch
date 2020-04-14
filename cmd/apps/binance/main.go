package main

import (
	"flag"
	"log"

	"github.com/Kifen/crypto-watch/pkg/binance"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatalf("Invalid number of arguments - found %d, requires [wsUrl] [sockfile]", len(flag.Args()))
	}

	wsUrl := flag.Args()[0]
	sockFile := flag.Args()[1]
	binance := binance.NewBinance(wsUrl, sockFile)
	err := binance.Serve()
	if err != nil {
		log.Fatalf("Error connecting appmanager client: %s", err)
	}
}
