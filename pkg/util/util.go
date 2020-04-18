package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Kifen/crypto-watch/pkg/proto"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/SkycoinProject/skycoin/src/util/logging"
)

func Logger(moduleName string) *logging.Logger {
	masterLogger := logging.NewMasterLogger()
	return masterLogger.PackageLogger(moduleName)
}

var log = Logger("util")
var once sync.Once

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
