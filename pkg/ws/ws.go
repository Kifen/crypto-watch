package ws

type Message struct {
	raw []byte
	t   int
}

type ExchangeConfig struct {
	WsUrl string
}

type Exchange interface {
	Connect() error
	WsRead() ([]byte, error)
	WsWrite(interface{}) error
	Stream()
	Run()
}

const (
	BinanceWsUrl = "wss://stream.binance.com:9443/ws/skybtc@ticker"
)
