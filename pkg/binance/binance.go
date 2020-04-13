package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"

	"github.com/Kifen/crypto-watch/pkg/ws"
)

type Binance struct {
	WsUrl    string
	sockFile string
	Dialer   *websocket.Dialer
	Header   http.Header
	WsConn   *websocket.Conn
	ErrCh    chan error
	AliveCh  chan struct{}
	DeadCh   chan struct{}
	logger   *logging.Logger
	mu       sync.Mutex
}

func NewBinance(url, sockFile string) ws.Exchange {
	return &Binance{
		sockFile: sockFile,
		WsUrl:    url,
		ErrCh:    make(chan error),
		DeadCh:   make(chan struct{}),
		logger:   util.Logger("Binance"),
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

func (b *Binance) Stream(symbol string) {
	symbol = strings.ToUpper(symbol)
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

func (b *Binance) Serve() error {
	listener, err := net.Listen("unix", b.sockFile)
	if err != nil {
		return err
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			b.logger.WithError(err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go b.handleConn(conn)
	}
}

func (b *Binance) handleConn(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != io.EOF {
			b.logger.Warnf("error on read: %s", err)
			break
		}

		subData, err := util.Deserialize(buf[:n])
		if err != nil {
			b.logger.Fatalf("Failed to deserialize data: %s", err)
		}

		b.Subscribe(subData)
	}
}

func (b *Binance) Subscribe(data ws.SubData) error {
	type subReq struct {
		method string   `json:"method"`
		params []string `json:"params"`
		id     int      `json:"id"`
	}

	s := fmt.Sprintf("%s@ticker", strings.ToLower(data.Symbol))
	req := &subReq{
		method: "SUBSCRIBE",
		params: []string{s},
		id:     data.Id,
	}
	msg, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe data: %s", err)
	}

	err = b.WsWrite(msg)
	if err != nil {
		return fmt.Errorf("failed to subscribe to stream: %s", err)
	}

	return nil
}
