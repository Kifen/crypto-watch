package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/bitly/go-simplejson"

	"github.com/Kifen/crypto-watch/pkg/proto"

	"github.com/SkycoinProject/skycoin/src/util/logging"
)

func Logger(moduleName string) *logging.Logger {
	masterLogger := logging.NewMasterLogger()
	return masterLogger.PackageLogger(moduleName)
}

var log = Logger("util")
var once sync.Once

const (
	URL     = "https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=%s"
	API_KEY = "5e1d9cde6e4b31856619d3061df490b1c62a4ea71b01da9dec4457e3a2c454ef"
)

func AppsDir(appsPath string) (string, error) {
	absPath, err := filepath.Abs(appsPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand path: %s", err)
	}

	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		return absPath, nil
	}

	if err := os.MkdirAll(absPath, 0750); err != nil {
		return "", fmt.Errorf("failed to create dir: %s", err)
	}

	return absPath, err
}

func GetBinaryPath(dir, name string) string {
	return filepath.Join(dir, name)
}

// FindConfigPath is used by a service to find a config file path in the following order:
// - From CLI argument.
// - From ENV.
// - From a list of default paths.
// If argsIndex < 0, searching from CLI arguments does not take place.
// Borrowed from github.com/SkycoinProject/skywire-mainnet/pkg/util/pathutil/configpath.go
func FindConfigPath(args []string, argsIndex int, env, defaultPath string) string {
	if argsIndex >= 0 && len(args) > argsIndex {
		path := args[argsIndex]
		log.Infof("using args[%d] as config path: %s", argsIndex, path)
		return path
	}

	if env != "" {
		if path, ok := os.LookupEnv(env); ok {
			log.Infof("using $%s as config path: %s", env, path)
			return path
		}
	}

	log.Debugf("config path is not explicitly specified, trying default path...")
	if _, err := os.Stat(defaultPath); err != nil {
		return defaultPath
	}

	log.Fatalf("config not found in defautl path %s", defaultPath)
	return ""
}

// UnlinkSocketFiles removes unix socketFiles from file system
func UnlinkSocketFiles(socketFiles ...string) error {
	for _, f := range socketFiles {
		if err := syscall.Unlink(f); err != nil {
			if !strings.Contains(err.Error(), "no such file or directory") {
				return err
			}
		}
	}

	return nil
}

func registerTypes() {
	gob.Register(proto.Symbol{})
	gob.Register(proto.AlertRes{})
	gob.Register(proto.AlertReq{})
}

func Serialize(data interface{}) ([]byte, error) {
	once.Do(func() {
		registerTypes()
	})

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&data)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func Deserialize(data []byte) (interface{}, error) {
	once.Do(func() {
		registerTypes()
	})

	var i interface{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func GetCryptoPrice(from, to string) (*float64, error) {
	finalUrl := fmt.Sprintf(URL, from, to)
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("authorization", API_KEY)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	j, err := simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}

	p, err := j.Get(strings.ToUpper(to)).Float64()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// QuickSort sorts an array using the quick sort algorithm
func QuickSort(arr []float32, left, right int) []float32 {
	if left < right {
		pivotPosition := partition(arr, left, right)
		QuickSort(arr, left, pivotPosition-1)
		QuickSort(arr, pivotPosition+1, right)
	}

	return arr
}

func partition(arr []float32, start, end int) int {
	swap := func(i, j int) {
		temp := arr[i]
		arr[i] = arr[j]
		arr[j] = temp
	}

	i, pivot := start-1, arr[end]
	for counter := start; counter <= end-1; counter++ {
		if arr[counter] <= pivot {
			i++
			swap(i, counter)
		}
	}

	newPivot := i + 1
	swap(newPivot, end)
	return newPivot
}
