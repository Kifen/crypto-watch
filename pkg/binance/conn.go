package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"

	"github.com/Kifen/crypto-watch/pkg/proto"
	"github.com/Kifen/crypto-watch/pkg/util"
)

func (b *Binance) handleWsConn(c *websocket.Conn, priceCh chan float32) error {
	for {
		msgByte, err := b.WsRead(c)
		if err != nil {
			return err
		}

		j, err := simplejson.NewJson(msgByte)
		if err != nil {
			return err
		}

		price := j.Get("c").Interface()
		b.logger.Info(price)
		p, err := strconv.ParseFloat(price.(string), 32)
		if err != nil {
			log.Fatal(err)
		}

		priceCh <- float32(p)
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
	err := b.initReqsSubscriptions()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", b.sockFile)
	if err != nil {
		return fmt.Errorf("failed to start binance server listener: %v", err)
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

		b.logger.Infof("New connection: %v", conn.LocalAddr())
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
		if err != nil {
			b.logger.Warnf("error on read: %v", err)
			return
		}

		reqData, err := util.Deserialize(buf[:n])
		if err != nil {
			b.logger.Fatalf("Failed to deserialize request data: %s", err)
		}

		switch v := reqData.(type) {
		case proto.AlertReq:
			b.logger.Info(v)
			err := b.registerRequest(&v)
			if err != nil {
				b.logger.Fatalf("Failed to register request: %v", err)
			}

			b.subscribe(strconv.Itoa(int(v.Id)), &v)

			b.logger.Infof("Request %d registered.", v.Id)
			go b.wsServe(&v, write)
		case proto.Symbol:
			b.validateSymbol(&v, write)
		}
	}
}

func (b *Binance) wsServe(req *proto.AlertReq, write func(b []byte, l *logging.Logger)) {
	endPoint := fmt.Sprintf("%s/%s@ticker", b.wsUrl, strings.ToLower(req.Req.Symbol))
	priceCh := make(chan float32)

	b.logger.Info("Establishing connection to binace ws...")
	c, _, err := websocket.DefaultDialer.Dial(endPoint, nil)
	if err != nil {
		b.logger.Fatalf("Failed to create a websocket connection to binance ws: %s", err)
	}

	b.logger.Info("Connection established to binance ws...")

	go b.handleWsConn(c, priceCh)
	pr := <-priceCh
	b.publish(pr)

	p := <-b.alertPriceCh

	if err := c.Close(); err != nil {
		b.logger.Fatalf("Failed to close websocket connection: %s", err)
	}

	res, err := util.Serialize(proto.AlertRes{
		Req:   req,
		Price: p,
	})

	if err != nil {
		log.Fatalf("Failed to serialize Response ReqData: %s", err)
	}

	write(res, b.logger)
	id := strconv.Itoa(int(req.Id))
	b.unsubscribe(id)
	b.logger.Infof("Unsubscribed request with id %s", id)
}
