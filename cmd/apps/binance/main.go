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
		log.Fatalf(
			"Invalid number of arguments - found %d, requires [wsUrl] [baseUrl] [sockfile] [redisUrl] [redisPassword]", len(flag.Args()))
	}

	wsUrl := flag.Args()[0]
	baseUrl := flag.Args()[1]
	sockFile := flag.Args()[2]
	redisUrl := flag.Args()[3]
	redisPassword := flag.Args()[4]

	binance, err := binance.NewBinance(sockFile, wsUrl, baseUrl, redisUrl, redisPassword)
	if err != nil {
		log.Fatalf("Failed to create binance app: %s", err)
	}

	err = binance.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
