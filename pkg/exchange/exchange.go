package exchange

import (
	"fmt"

	pb "github.com/Kifen/crypto-watch/pkg/proto"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"
)

type AppConfig struct {
	Exchange   string `json:"exchange"`
	WsUrl      string `json:"ws_url"`
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

	return &Exchange{
		Store:      s,
		Logger:     util.Logger("Exchange"),
		Srv:        NewServer(),
		AppManager: appManager,
	}, nil
}

func (e *Exchange) ManageServerConn() {
	var fn = func(req *pb.ExchangeReq) {
		err := e.AppManager.appExists(req.Exchange)
		if err == ErrAppNotFound {
			e.Srv.ResCH <- &pb.ExchangeRes{
				Req:     req,
				Message: fmt.Sprint("Exchange %s is not supported", req.Exchange)}
			return
		}

		if alive := e.AppManager.appIsAlive(req.Exchange); !alive {
			if err := e.AppManager.StartApp(req.Exchange); err != nil {
				e.Logger.Warnf("Failed to start %s app: %s", req.Exchange, err)
				return
			}
		}

		data := ReqData{
			Symbol: req.Req.Symbol,
			Id:     int(req.Id),
		}
		if err := e.AppManager.SendData(req.Exchange, data); err != nil {
			e.Logger.Errorf("Failed to send data to %s server: %s", req.Exchange, err)
			return
		}
	}

	go e.Srv.handleConn(fn)
}
