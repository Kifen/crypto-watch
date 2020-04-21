package main

import (
	"flag"

	"github.com/Kifen/crypto-watch/pkg/binance"
	"github.com/Kifen/crypto-watch/pkg/util"
)

func main() {
	log := util.Logger("Binance")
	flag.Parse()

	if len(flag.Args()) < 3 {
		log.Fatalf("Invalid number of arguments - found %d, requires [wsUrl] [baseUrl] [sockfile]", len(flag.Args()))
	}

	wsUrl := flag.Args()[0]
	baseUrl := flag.Args()[1]
	sockFile := flag.Args()[2]
	binance := binance.NewBinance(sockFile, wsUrl, baseUrl)

	err := binance.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
