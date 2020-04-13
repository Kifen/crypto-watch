package binance

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"

	"github.com/Kifen/crypto-watch/pkg/ws"
)

type Binance struct {
	WsUrl   string
	Dialer  *websocket.Dialer
	Header  http.Header
	WsConn  *websocket.Conn
	ErrCh   chan error
	AliveCh chan struct{}
	DeadCh  chan struct{}
	logger  *logging.Logger
}

func NewBinance(url string) ws.Exchange {
	return &Binance{
		WsUrl:  url,
		ErrCh:  make(chan error),
		DeadCh: make(chan struct{}),
		logger: util.Logger("Binance"),
	}
}

func (b *Binance) Connect() error {
	conn, _, err := b.Dialer.Dial(b.WsUrl, b.Header)
	if err != nil {
		return err
	}

	b.WsConn = conn
	return nil
}

func (b *Binance) WsRead() ([]byte, error) {
	_, p, err := b.WsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (b *Binance) WsWrite(data interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return b.WsConn.WriteMessage(websocket.TextMessage, msg)
}

func (b *Binance) Stream() {
	for {
		msgByte, err := b.WsRead()
		if err != nil {
			log.Fatal(err)
		}

		//log.Println("Logging data:\n", string(msgByte))
		j, err := simplejson.NewJson(msgByte)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("#PRICE SKYBTC: ", j.Get("c"))
		os.Exit(1)
	}
}

func (b *Binance) Run() {
	err := b.Connect()
	if err != nil {
		b.ErrCh <- err
		return
	}

	b.logger.Info("Binance Initialized...")
	go b.Stream()
}
