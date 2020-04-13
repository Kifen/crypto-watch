package main

import (
	"flag"
	"log"

	"github.com/Kifen/crypto-watch/pkg/binance"
)

func main() {
	//url := "wss://stream.binance.com:9443/ws/skybtc@ticker"
	flag.Parse()
	if len(flag.Args()) <3 {
		log.Fatalf("Invalid number of arguments - found %d, requires [wsUrl] [symbol] [sockfile]", len(flag.Args()))
	}

	wsUrl := flag.Args()[0]
	symbol := flag.Args()[1]
	sockFile := flag.Args()[2]
	binance := binance.NewBinance(wsUrl, sockFile)
	err := binance.Connect()
	if err != nil {
		log.Fatalf("Error connecting appmanager client: %s", err)
	}

	log.Println("Connected to binance stream ws...")
	binance.Stream(symbol)
}
