package exchange

import (
	"errors"
	"fmt"
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

func (a *AppManager) StartApp(exchangeName string) error {
	exchange := strings.ToLower(exchangeName)

	if _, ok := a.apps[exchange]; ok {
		if _, alive := a.appsState[exchange]; !alive {
			err := a.start(a.appsPath, exchange)
			if err != nil {
				return err
			}

			a.appMu.Lock()
			a.appsState[exchange] = struct{}{}
			a.appMu.Unlock()

			return nil
		}

		return ErrAppAlreadyStarted
	}

	return ErrAppNotFound
}

func (a *AppManager) start(dir, name string) error {
	wsUrl := a.apps[name].WsUrl
	binaryPath := util.GetBinaryPath(dir, name)
	cmd := exec.Command(binaryPath, wsUrl)

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
