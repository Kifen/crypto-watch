package util

import (
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

func Logger(moduleName string) *logging.Logger {
	masterLogger := logging.NewMasterLogger()
	return masterLogger.PackageLogger(moduleName)
}
