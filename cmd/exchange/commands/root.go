package commands

import (
	"github.com/spf13/cobra"
	"log"

	"github.com/Kifen/crypto-watch/pkg/exchange"
)

type runCfg struct {
	redisAddr     string
	redisPassword string
	serverAddr    string
}

var cfg *runCfg

var rootCmd = &cobra.Command{
	Use:   "crypto-watch",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		exchange := initExchange(cfg.redisAddr, cfg.redisPassword)
		exchange.Srv.StartServer(cfg.serverAddr)
	},
}

func init() {
	cfg = &runCfg{}
	rootCmd.Flags().StringVarP(&cfg.redisAddr, "redis-url", "", "redis://localhost:6379", "redis address")
	rootCmd.Flags().StringVarP(&cfg.redisPassword, "password", "p", "none", "redis password")
	rootCmd.Flags().StringVarP(&cfg.serverAddr, "server-addr", "", ":2000", "grpc server address")
}

func initExchange(url, password string) *exchange.Exchange {
	ex, err := exchange.NewExchange(url, password)
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
