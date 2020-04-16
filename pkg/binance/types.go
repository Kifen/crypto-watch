package binance

type Symbol struct {
	symbol string `json:"Symbol"`
}

type Response struct {
	symbols []Symbol `json:"symbols"`
}
