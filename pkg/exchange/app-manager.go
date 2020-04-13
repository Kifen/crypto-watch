package exchange

import (
	"errors"
	"fmt"
	"github.com/Kifen/crypto-watch/pkg/ws"
	"net"
	"os/exec"
	"strings"
	"sync"

	"github.com/Kifen/crypto-watch/pkg/util"
	"github.com/SkycoinProject/skycoin/src/util/logging"
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
	appMu sync.Mutex
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

func (a *AppManager) StartApp(exchangeName, symbol string, id int) error {
	exchange := strings.ToLower(exchangeName)

	if _, ok := a.apps[exchange]; ok {
		if _, alive := a.appsState[exchange]; !alive {
			err := a.start(a.appsPath, exchange, symbol)
			if err != nil {
				return err
			}

			a.appMu.Lock()
			a.appsState[exchange] = struct{}{}
			a.appMu.Unlock()

			return nil
		} else {
			sock := a.apps[exchangeName].SocketFile
			subData := ws.SubData{
				Symbol: symbol,
				Id: id,
			}

			data, err := util.Serialize(subData)
			if err != nil {
				a.log.Fatalf("failed to serialize data: %s", err)
			}

			if err := a.sendData(data, sock); err != nil {
				a.log.Errorf("Failed to send symbol %s to %s server.", symbol, exchange)
			}
		}

		return ErrAppAlreadyStarted
	}

	return ErrAppNotFound
}

func (a *AppManager) start(dir, name, symbol string) error {
	wsUrl := a.apps[name].WsUrl
	binaryPath := util.GetBinaryPath(dir, name)
	cmd := exec.Command(binaryPath, wsUrl, symbol)

	if err := cmd.Start(); err != nil {
		return err
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

