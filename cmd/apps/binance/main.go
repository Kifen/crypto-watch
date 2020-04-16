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
	baseUrl := flag.Args()[1]
	sockFile := flag.Args()[2]
	binance := binance.NewBinance(sockFile, wsUrl, baseUrl)
	err := binance.Serve()
	if err != nil {
		log.Fatalf("Error connecting appmanager client: %s", err)
	}
}
