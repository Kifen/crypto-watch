package exchange

import (
	"fmt"
	pb "github.com/Kifen/crypto-watch/pkg/proto"

	"github.com/Kifen/crypto-watch/pkg/util"
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

type AppConfig struct {
	Exchange string `json:"exchange"`
	WsUrl    string `json:"ws_url"`
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

	return &Exchange{
		Store:      s,
		Logger:     util.Logger("Exchange"),
		Srv:        NewServer(),
		AppManager: appManager,
	}, nil
}

func (e *Exchange) ManageServerConn() {
	var fn = func(req *pb.ExchangeReq) {
		if err := e.AppManager.StartApp(req.Exchange); err != nil && err != ErrAppAlreadyStarted {
			if err == ErrAppNotFound {
				e.Srv.ResCH <- &pb.ExchangeRes{
					Req:     req,
					Message: fmt.Sprint("Exchange %s is not supported", req.Exchange)}
			}
			e.Srv.ResCH <- &pb.ExchangeRes{
				Req:     req,
				Message: fmt.Sprint("Exchange service %s is currently down.", req.Exchange)}
		}
	}

	go e.Srv.handleConn(fn)
}
