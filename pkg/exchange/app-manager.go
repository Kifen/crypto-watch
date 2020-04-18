package exchange

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

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
	appMu     sync.Mutex
	once      sync.Once
	msgCh     chan interface{}
}

func NewAppManager(appsPath string, exchangeApps []AppConfig) (*AppManager, error) {
	path, err := util.AppsDir(appsPath)
	if err != nil {
		return nil, fmt.Errorf("Invalid apps path: %s", err)
	}
	return &AppManager{
		log:       util.Logger("appmanager-rpc"),
		apps:      makeApps(exchangeApps),
		appsState: map[string]struct{}{},
		appsPath:  path,
		msgCh:     make(chan interface{}, 100),
	}, nil
}

func makeApps(exchangeApps []AppConfig) map[string]AppConfig {
	apps := make(map[string]AppConfig)
	for _, app := range exchangeApps {
		apps[app.Exchange] = app
	}

	return apps
}

func (a *AppManager) appExists(exchange string) error {
	if _, ok := a.apps[exchange]; !ok {
		return ErrAppNotFound
	}

	return nil
}

func (a *AppManager) appIsAlive(exchange string) bool {
	_, alive := a.appsState[exchange]
	return alive
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

func (a *AppManager) SendData(exchange string, req interface{}) error {
	if err := a.appExists(exchange); err != nil {
		return err
	}

	if alive := a.appIsAlive(exchange); !alive {
		if err := a.StartApp(exchange); err != nil {
			return err
		}
		a.log.Infof("Successfully started %s server.", exchange)
	}

	data, err := util.Serialize(req)
	if err != nil {
		a.log.Fatalf("failed to serialize data: %s", err)
		return err
	}

	if err := a.sendData(data, a.apps[exchange].SocketFile); err != nil {
		return err
	}

	return nil
}

func (a *AppManager) start(dir, name, sockFile string) error {
	if err := util.UnlinkSocketFiles(sockFile); err != nil {
		return fmt.Errorf("failed to remove unix socket file: %s", err)
	}

	wsUrl := a.apps[name].WsUrl
	baseUrl := a.apps[name].BaseUrl
	binaryPath := util.GetBinaryPath(dir, name)

	//cmd := exec.Command(binaryPath, wsUrl, baseUrl, sockFile)
	cmd := &exec.Cmd{
		Path:   binaryPath,
		Args:   []string{binaryPath, wsUrl, baseUrl, sockFile},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start binance app: %s", err)
	}

	time.Sleep(5 * time.Second)
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

	a.log.Info("Writing request to app server.")
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	a.log.Infof("Wrote %d bytes", n)

	a.once.Do(func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				a.log.Fatalf("Error on read: %s", err)
			}

			resData, err := util.Deserialize(buf[:n])
			if err != nil {
				a.log.Fatalf("Failed to deserialize response data: %s", err)
			}

			a.msgCh <- resData
		}
	})
	return nil
}
