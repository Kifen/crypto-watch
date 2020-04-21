package exchange

import (
	"fmt"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	pb "github.com/Kifen/crypto-watch/pkg/proto"

	"github.com/Kifen/crypto-watch/pkg/util"
)

type AppConfig struct {
	Exchange   string `json:"exchange"`
	WsUrl      string `json:"ws_url"`
	BaseUrl    string `json:"base_url"`
	SocketFile string `json:"socket_file"`
}

type Config struct {
	AppsPath      string      `json:"apps_path"`
	RedisAddr     string      `json:"redis_addr"`
	RedisPassword string      `json:"redis_password"`
	ServerAddr    string      `json:"server_addr"`
	ExchangeApps  []AppConfig `json:"exchange_apps"`
}

type Exchange struct {
	Store      *RedisStore
	Logger     *logging.Logger
	Srv        *Server
	AppManager *AppManager
}

func NewExchange(redisUrl, password, appsPath string, exchangeApps []AppConfig) (*Exchange, error) {
	s, err := NewRedisStore(redisUrl, password)
	if err != nil {
		return nil, err
	}

	// setup the app manager.
	appManager, err := NewAppManager(appsPath, exchangeApps)
	if err != nil {
		return nil, fmt.Errorf("failed to setup app manager: %s", err)
	}

	var supportedExchanges = func(exchangeName string) bool {
		_, ok := appManager.apps[exchangeName]
		return ok
	}

	return &Exchange{
		Store:      s,
		Logger:     util.Logger("Exchange"),
		Srv:        NewServer(supportedExchanges),
		AppManager: appManager,
	}, nil
}

func (e *Exchange) ManageServerConn() {
	e.Logger.Info("Managing connection...")
	for {
		select {
		case req := <-e.Srv.alertReqCh:
			e.handlePriceRequest(req)
		case req := <-e.Srv.symbolReqCH:
			e.handleSymbolRequest(req)
		case res := <-e.AppManager.msgCh:
			switch v := res.(type) {
			case pb.Symbol:
				e.Srv.symbolResCH <- &v
			case pb.AlertRes:
				e.Srv.alertResCh <- &v
			}
		}
	}
}

func (e *Exchange) handlePriceRequest(req *pb.AlertReq) {
	if err := e.AppManager.SendData(req.ExchangeName, *req); err != nil {
		e.Logger.Errorf("Failed to send data to %s server: %s", req.ExchangeName, err)
		e.Srv.alertResCh <- &pb.AlertRes{
			Req:     req,
			Price:   0.00,
			Message: err.Error(),
		}
	}
}

func (e *Exchange) handleSymbolRequest(req *pb.Symbol) {
	if err := e.AppManager.SendData(req.ExchangeName, *req); err != nil {
		e.Logger.Errorf("Failed to send data to %s server: %s", req.ExchangeName, err)
		e.Srv.symbolResCH <- &pb.Symbol{
			Id:           req.Id,
			ExchangeName: req.ExchangeName,
			Symbol:       req.Symbol,
			Valid:        false,
			Message:      err.Error(),
		}
	}
}
