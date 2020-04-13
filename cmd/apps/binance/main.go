package main

import (
	"log"

	"github.com/Kifen/crypto-watch/pkg/binance"
)

func main() {
	url := "wss://stream.binance.com:9443/ws/skybtc@ticker"
	binance := binance.NewBinance(url)
	err := binance.Connect()
	if err != nil {
		log.Fatalf("Error connecting appmanager client: %s", err)
	}

	log.Println("Connected to binance stream ws...")
	binance.Stream()
}
