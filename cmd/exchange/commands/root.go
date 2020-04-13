package commands

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/Kifen/crypto-watch/pkg/util"

	"github.com/spf13/cobra"

	"github.com/Kifen/crypto-watch/pkg/exchange"
)

const configEnv = "CW_CONFIG"

var log = util.Logger("crypro-watch")

type runCfg struct {
	configPath   *string
	args         []string
	cfgFromStdin bool
	conf         exchange.Config
	logger       *logging.Logger
}

var cfg *runCfg

var rootCmd = &cobra.Command{
	Use:   "exchange [config-path]",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.args = args
		cfg.logger = util.Logger("")
		readConfig()
		exchange := initExchange(cfg.conf.RedisAddr, cfg.conf.RedisPassword, cfg.conf.AppsPath, cfg.conf.ExchangeApps)
		go exchange.Srv.StartServer(cfg.conf.ServerAddr)
		go exchange.ManageServerConn()
		waitOsSignals()
	},
}

func init() {
	cfg = &runCfg{}
	//rootCmd.Flags().StringVarP(&cfg.redisAddr, "redis-url", "", "redis://localhost:6379", "redis address")
	//rootCmd.Flags().StringVarP(&cfg.redisPassword, "password", "p", "none", "redis password")
	//rootCmd.Flags().StringVarP(&cfg.serverAddr, "server-addr", "", ":2000", "grpc server address")
	rootCmd.Flags().BoolVarP(&cfg.cfgFromStdin, "stdin", "i", false, "read config from STDIN")
	//rootCmd.Flags().StringVarP(&cfg.appsPath, "apppath", "", "./bin/apps", "path to exchange apps binaries")
}

func initExchange(url, password, appsPath string, exchangeApps []exchange.AppConfig) *exchange.Exchange {
	ex, err := exchange.NewExchange(url, password, appsPath, exchangeApps)
	if err != nil {
		log.Fatalf("failed to setup crypto watch service: %s", err)
	}

	return ex
}

// Execute executes root CLI command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func readConfig() {
	var rdr io.Reader

	if !cfg.cfgFromStdin {
		configPath := util.FindConfigPath(cfg.args, 0, configEnv, defaultConfigPath())
		file, err := os.Open(filepath.Clean(configPath))
		if err != nil {
			log.Fatalf("Failed to open config: %s", err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Warnf("Failed to close config file: %v", err)
			}
		}()

		log.Infof("Reading config from %v", configPath)

		rdr = file
		cfg.configPath = &configPath
	} else {
		log.Info("Reading config from STDIN")
		rdr = bufio.NewReader(os.Stdin)
	}

	cfg.conf = exchange.Config{}
	if err := json.NewDecoder(rdr).Decode(&cfg.conf); err != nil {
		log.Fatalf("Failed to decode %s: %s", rdr, err)
	}
	log.Info(cfg.conf)
}

func defaultConfigPath() (path string) {
	if wd, err := os.Getwd(); err != nil {
		path = filepath.Join(wd, "exchange.json")
	}

	return
}

func waitOsSignals() {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}...)
	<-ch

	go func() {
		select {
		case s := <-ch:
			cfg.logger.Fatalf("Received signal %s: terminating", s)
		}
	}()
}
