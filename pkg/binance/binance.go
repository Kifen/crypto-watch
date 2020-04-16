package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Kifen/crypto-watch/pkg/proto"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

type Binance struct {
	sockFile string
	logger   *logging.Logger
	wsUrl    string
	BaseUrl  string
	ErrCh    chan error
	wg       sync.WaitGroup
}

func NewBinance(sockFile, wsUrl, baseUrl string) *Binance {
	return &Binance{
		sockFile: sockFile,
		wsUrl:    wsUrl,
		BaseUrl:  baseUrl,
		logger:   util.Logger("Binance"),
		ErrCh:    make(chan error),
	}
}

func (b *Binance) handleWsConn(c *websocket.Conn, priceCh chan float32) error {
	b.logger.Info("Connection established to binance ws...")

	for {
		msgByte, err := b.WsRead(c)
		if err != nil {
			return err
		}

		log.Println("Logging data:\n", string(msgByte))
		j, err := simplejson.NewJson(msgByte)
		if err != nil {
			return err
		}

		price := j.Get("c").Interface()
		priceCh <- price.(float32)
		//log.Println("#PRICE SKYBTC: ", j.Get("c"))
		//os.Exit(1)
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

		b.handleServerConn(conn)
	}

	return nil
}

func (b *Binance) handleServerConn(conn net.Conn) {
	var write = func(b []byte, logger *logging.Logger) {
		n, err := conn.Write(b)
		if err != nil {
			logger.Fatalf("Failed to write Response to unix client")
		}

		logger.Infof("Wrote %d bytes to unix client", n)
	}

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

		switch v := reqData.(type) {
		case proto.AlertReq:
			go b.wsServe(&v, write)
		case proto.Symbol:
			b.validateSymbol(&v, write)
		}
	}
}

func (b *Binance) wsServe(req *proto.AlertReq, write func(b []byte, l *logging.Logger)) {
	endPoint := fmt.Sprintf("%s/%s@ticker", b.wsUrl, strings.ToLower(req.Req.Symbol))
	priceCh := make(chan float32)
	alertPriceCh := make(chan float32)

	b.logger.Info("Establishing connection to binace ws...")
	c, _, err := websocket.DefaultDialer.Dial(endPoint, nil)
	if err != nil {
		b.logger.Fatalf("Failed to create a websocket connection to binance ws: %s", err)
	}

	b.logger.Info("Connection established to binance ws...")

	go b.handleWsConn(c, priceCh)
	go b.alert(req.Req.Action, req.Req.Price, priceCh, alertPriceCh)

	p := <-alertPriceCh
	if err := c.Close(); err != nil {
		b.logger.Fatalf("Failed to close websocket connection: %s", err)
	}

	res, err := util.Serialize(proto.AlertRes{
		Req:   req,
		Price: p,
	})

	if err != nil {
		log.Fatalf("Failed to serialize Response req: %s", err)
	}

	write(res, b.logger)
}

func (b *Binance) alert(action string, price float32, priceCh, alertPriceCh chan float32) {
	switch action {
	case "gt":
		for {
			select {
			case rePrice := <-priceCh:
				if rePrice > price {
					alertPriceCh <- rePrice
					return
				}
			}
		}
	case "lt":
		for {
			select {
			case rePrice := <-priceCh:
				if rePrice < price {
					alertPriceCh <- rePrice
					return
				}
			}
		}
	}
}

func (b *Binance) validateSymbol(s *proto.Symbol, fn func(b []byte, l *logging.Logger)) {
	var find = func(slice []Symbol, val string) bool {
		for _, item := range slice {
			if item.symbol == val {
				return true
			}
		}
		return false
	}

	r, err := b.getExchangeSymbols()
	if err != nil {
		b.logger.Fatalf("Failed to get exchange symbols: %s", err)
	}

	isValid := find(r.symbols, s.Symbol)
	res, err := util.Serialize(proto.Symbol{
		ExchangeName: s.ExchangeName,
		Symbol:       s.Symbol,
		Valid:        isValid,
	})

	if err != nil {
		log.Fatalf("Failed to serialize Response req: %s", err)
	}

	fn(res, b.logger)
}

func (b *Binance) getExchangeSymbols() (*Response, error) {
	endpoint := "api/v3/exchangeInfo"
	url := fmt.Sprintf("%s/%s", b.BaseUrl, endpoint)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Response body: %s", err)
	}

	var symbols *Response
	err = json.Unmarshal(body, &symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall data: %s", err)
	}

	b.logger.Infof("Unmarshalled data: %s", symbols)

	return symbols, nil
}
