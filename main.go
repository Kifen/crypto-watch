package main

import (
	"github.com/Kifen/crypto-watch/pkg/ws"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	url := "wss://stream.binance.com:9443/ws/skybtc@ticker"
	conn, err := ws.Connect(url, &websocket.Dialer{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to binance stream ws...")
	ws.Streamticker(conn)
}
