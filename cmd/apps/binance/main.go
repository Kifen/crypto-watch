package main

import (
	"flag"
	"log"

	"github.com/Kifen/crypto-watch/pkg/binance"
)

func main() {
	//url := "wss://stream.binance.com:9443/ws/skybtc@ticker"
	flag.Parse()
	if len(flag.Args()) <1 {
		log.Fatalf("Invalid number of arguments - found %d, requires [wsUrl", len(flag.Args()))
	}

	wsUrl := flag.Args()[0]
	binance := binance.NewBinance(wsUrl)
	err := binance.Connect()
	if err != nil {
		log.Fatalf("Error connecting appmanager client: %s", err)
	}

	log.Println("Connected to binance stream ws...")
	binance.Stream()
}
