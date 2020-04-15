package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

type Binance struct {
	sockFile string
	logger   *logging.Logger
	wsUrl    string
	ErrCh    chan error
	wg       sync.WaitGroup
}

func NewBinance(sockFile, url string) *Binance {
	return &Binance{
		sockFile: sockFile,
		wsUrl:    url,
		logger:   util.Logger("Binance"),
		ErrCh:    make(chan error),
	}
}

func (b *Binance) WsServe(wsUrl string, priceCh chan float64) error {
	b.logger.Info("Establishing connection to binace ws...")
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}

	b.logger.Info("Connection established to binance ws...")

	for {
		msgByte, err := b.WsRead(conn)
		if err != nil {
			return err
		}

		log.Println("Logging data:\n", string(msgByte))
		j, err := simplejson.NewJson(msgByte)
		if err != nil {
			return err
		}

		log.Println("#PRICE SKYBTC: ", j.Get("c"))
		os.Exit(1)
	}
}

func (b *Binance) WsRead(conn *websocket.Conn) ([]byte, error) {
	_, p, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (b *Binance) WsWrite(conn *websocket.Conn, data interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, msg)
}

func (b *Binance) Serve() error {
	listener, err := net.Listen("unix", b.sockFile)
	if err != nil {
		return err
	}

	b.logger.Infof("Binance server listening on unix socket: %s", b.sockFile)

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

		go b.handleServerConn(conn)
	}

	return nil
}

func (b *Binance) handleServerConn(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != io.EOF {
			b.logger.Warnf("error on read: %s", err)
			break
		}

		reqData, err := util.Deserialize(buf[:n])
		if err != nil {
			b.logger.Fatalf("Failed to deserialize request data: %s", err)
		}

		b.serve(reqData, conn)
	}
}

func (b *Binance) serve(data util.ReqData, conn net.Conn) {
	endPoint := fmt.Sprintf("%s/%s@ticker", b.wsUrl, strings.ToLower(data.Symbol))
	go func() {
		priceCh := make(chan float64)
		b.WsServe(endPoint, priceCh)
		price := <-priceCh

		b, err := util.Serialize(util.ResData{
			Symbol: data.Symbol,
			Id:     data.Id,
			Price:  price,
		})

		if err != nil {
			log.Fatalf("Failed to serialize response data: %s", err)
		}

		// write price data to client
		n, err := conn.Write(b)
		if err != nil {
			log.Fatalf("Failed to write data to unix client")
		}

		log.Printf("Wrote %d bytes to unix client", n)
	}()
}
