package exchange

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"
)

var (
	//ErrAppNotFound is returned when an app is not found.
	ErrAppNotFound = errors.New("app not found")

	//ErrAppNotFound is returned when an app is not found.
	ErrAppAlreadyStarted = errors.New("app is already started")
)

type AppManager struct {
	log       *logging.Logger
	apps      map[string]AppConfig
	appsState map[string]struct{}
	appsPath  string
	appMu     sync.Mutex
}

func NewAppManager(appsPath string, exchangeApps []AppConfig) (*AppManager, error) {
	path, err := util.AppsDir(appsPath)
	if err != nil {
		return nil, fmt.Errorf("Invalid apps path: %s", err)
	}
	return &AppManager{
		log:      util.Logger("appmanager-rpc"),
		apps:     makeApps(exchangeApps),
		appsPath: path,
	}, nil
}

func makeApps(exchangeApps []AppConfig) map[string]AppConfig {
	apps := make(map[string]AppConfig)
	for _, app := range exchangeApps {
		apps[app.Exchange] = app
	}

	return apps
}

/*func StartAppServer(r *AppManager, addr string) {
	rpcS := rpc.NewServer()
	rpcG := &AppManager{
		log: util.Logger("appmanager-rpc"),
	}

	if err := rpcS.Register(rpcG); err != nil {
		r.log.Fatalf("failed to register appmanager-rpc server: %s", err)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		r.log.Fatalf("failed to set appmanager-rpc listener: %s", err)
	}

	r.log.Info("App-rpc server up...")
	rpc.Accept(lis)
}*/

func (a *AppManager) appExists(exchange string) error {
	if _, ok := a.apps[exchange]; !ok {
		return ErrAppNotFound
	}

	return nil
}

func (a *AppManager) appIsAlive(exchange string) bool {
	if _, alive := a.appsState[exchange]; alive {
		return alive
	}

	return false
}

func (a *AppManager) StartApp(exchangeName string) error {
	exchange := strings.ToLower(exchangeName)
	sockFile := a.apps[exchangeName].SocketFile

	err := a.start(a.appsPath, exchange, sockFile)
	if err != nil {
		return err
	}

	a.appMu.Lock()
	a.appsState[exchange] = struct{}{}
	a.appMu.Unlock()

	return nil
}

func (a *AppManager) SendData(exchange string, subData ReqData) error {
	data, err := util.Serialize(subData)
	if err != nil {
		a.log.Fatalf("failed to serialize data: %s", err)
		return err
	}

	if err := a.sendData(data, a.apps[exchange].SocketFile); err != nil {
		a.log.Errorf("Failed to send symbol %s to %s server.", subData.Symbol, exchange)
		return err
	}

	return nil
}

func (a *AppManager) start(dir, name, sockFile string) error {
	wsUrl := a.apps[name].WsUrl
	binaryPath := util.GetBinaryPath(dir, name)
	cmd := exec.Command(binaryPath, wsUrl, sockFile)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start binance app: %s", err)
	}

	return nil
}

func (a *AppManager) sendData(data []byte, socketFile string) error {
	conn, err := net.Dial("unix", socketFile)
	if err != nil {
		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			a.log.Errorf("error closing client: %v", err)
		}
	}()

	n, err := conn.Write(data)
	if err != nil {
		return err
	}

	a.log.Infof("Wrote %d bytes", n)
	return nil
}
