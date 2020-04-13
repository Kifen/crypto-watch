package ws

type Message struct {
	raw []byte
	t   int
}

type ExchangeConfig struct {
	WsUrl string
}

type SubData struct {
	Symbol string
	Id     int
}

type Exchange interface {
	Connect() error
	WsRead() ([]byte, error)
	WsWrite(interface{}) error
	Stream(string)
}
