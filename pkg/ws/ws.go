package ws

type Message struct {
	raw []byte
	t   int
}

type ExchangeConfig struct {
	WsUrl string
}

type ReqData struct {
	Symbol string
	Id     int
}

type ResData struct {
	Symbol string
	Id     int
	Price  float64
}

type Exchange interface {
	Connect() error
	WsRead() ([]byte, error)
	WsWrite(interface{}) error
	Stream(ReqData)
}
