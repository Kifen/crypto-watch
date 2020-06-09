package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/bitly/go-simplejson"
	"github.com/go-redis/redis"

	"github.com/Kifen/crypto-watch/pkg/proto"
	"github.com/Kifen/crypto-watch/pkg/store"
	"github.com/Kifen/crypto-watch/pkg/util"
)

const (
	GT = "gt"
	LT = "lt"
)

var (
	// ErrAlreadyRegistered indicates that request ID is already in use.
	ErrAlreadyRegistered = errors.New("ID already registered")

	// ErrRequestNotFound indicates that requeste is not registered.
	ErrRequestNotFound = errors.New("request not found")

	// ErrKeysNotFound indicates that no keys registered.
	ErrKeysNotFound = errors.New("keys not found")
)

type Binance struct {
	sockFile     string
	logger       *logging.Logger
	wsUrl        string
	BaseUrl      string
	ErrCh        chan error
	alertPriceCh chan float32
	wg           sync.WaitGroup
	gtArr        []float32
	ltArr        []float32
	reqChs       map[string]proto.AlertReq
	redisStore   *store.RedisStore
	rm           sync.RWMutex
}

func NewBinance(sockFile, wsUrl, baseUrl, redisUrl, redisPassword string) (*Binance, error) {
	s, err := store.NewRedisStore(redisUrl, redisPassword)
	if err != nil {
		return nil, err
	}

	return &Binance{
		sockFile:     sockFile,
		wsUrl:        wsUrl,
		BaseUrl:      baseUrl,
		logger:       util.Logger("Binance"),
		ErrCh:        make(chan error),
		alertPriceCh: make(chan float32),
		redisStore:   s,
		reqChs:       make(map[string]proto.AlertReq),
	}, nil
}

func (b *Binance) subscribe(id string, req proto.AlertReq) {
	b.rm.Lock()
	if _, ok := b.reqChs[id]; !ok {
		b.reqChs[id] = req
	}
	b.rm.Unlock()
}

func (b *Binance) unsubscribe(id string) {
	delete(b.reqChs, id)
}

func (b *Binance) publish(price float32) {
	b.rm.Lock()
	for _, v := range b.reqChs {
		go b.check(v.Req.Action, v.Req.Price, price, v.ExchangeName, v.Req.Symbol)
	}
	b.rm.Unlock()
}

func (b *Binance) check(action string, price, rePrice float32, exchange, symbol string) {
	switch action {
	case GT:
		from, to := symbol[3:], "USD"
		p, err := b.calcPrice(from, to, float64(rePrice))
		if err != nil {
			return
		}

		if float32(p) > price {
			b.logger.Infof("Alert: %s goes over %f on %s. Current price is %f USD", strings.ToUpper(symbol[0:3]), price, exchange, p)
			b.alertPriceCh <- float32(p)
			return
		}
	case LT:
		from, to := symbol[3:], "USD"
		p, err := b.calcPrice(from, to, float64(rePrice))
		if err != nil {
			return
		}

		if float32(p) < price {
			b.logger.Infof("Alert: %s goes below %f on %s. Current price is %f USD", strings.ToUpper(symbol[0:3]), price, exchange, p)
			b.alertPriceCh <- float32(p)
			return
		}
	}
}

func (b *Binance) validateSymbol(s *proto.Symbol, fn func(b []byte, l *logging.Logger)) {
	v, err := b.getSymbolsAndValidate(s.Symbol)
	if err != nil {
		b.logger.Fatalf("Failed to validate symbol: %s", err)
	}

	b.logger.Infof("Symbol %s is valid on %s exchange.", s.Symbol, s.ExchangeName)

	res, err := util.Serialize(proto.Symbol{
		Id:           s.Id,
		ExchangeName: s.ExchangeName,
		Symbol:       s.Symbol,
		Valid:        v,
	})

	if err != nil {
		log.Fatalf("Failed to serialize Response ReqData: %s", err)
	}

	fn(res, b.logger)
}

func (b *Binance) getSymbolsAndValidate(symbol string) (bool, error) {
	endpoint := "api/v3/exchangeInfo"
	url := fmt.Sprintf("%s/%s", b.BaseUrl, endpoint)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to get resource: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read Response body: %s", err)
	}

	j, err := simplejson.NewJson(body)
	if err != nil {
		return false, err
	}

	symbolLen := len(j.Get("symbols").MustArray())
	b.logger.Infof("Validating symbol %s.", symbol)

	for i := 0; i < symbolLen; i++ {
		s := j.Get("symbols").GetIndex(i).Get("symbol").Interface()
		if strings.ToLower(s.(string)) == symbol {
			return true, nil
		}
	}

	return false, fmt.Errorf("symbol %s is not supported on exchange.", symbol)
}

func (b *Binance) calcPrice(from, to string, rePrice float64) (float64, error) {
	p, err := util.GetCryptoPrice(strings.ToUpper(from), to)
	if err != nil {
		return 0, err
	}

	v := rePrice * (*p)
	return v, nil
}

func (b *Binance) initBinance() ([]float32, []float32, error) {
	var (
		tempGt, tempLt []float32
	)

	gts, err := b.getRequestsById(GT)
	if err != nil && err != ErrKeysNotFound {
		return nil, nil, err
	}

	for _, req := range gts {
		tempGt = append(tempGt, req.Req.Price)
	}

	lts, err := b.getRequestsById(LT)
	if err != nil && err != ErrKeysNotFound {
		return nil, nil, err
	}

	for _, req := range lts {
		tempLt = append(tempLt, req.Req.Price)
	}

	tempGt = util.QuickSort(tempGt, 0, len(tempGt)-1)
	tempLt = util.QuickSort(tempLt, 0, len(tempLt)-1)

	return tempGt, tempLt, nil
}

func (b *Binance) getRequestsById(key string) ([]*proto.AlertReq, error) {
	reqPrices, err := b.redisStore.Client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	var keys []string

	for _, reqPrice := range reqPrices {
		v, err := strconv.ParseFloat(reqPrice, 32)
		if err != nil {
			return nil, fmt.Errorf("error converting string to float: %v", err)
		}
		keys = append(keys, reqKey(float32(v)))
	}

	if len(keys) == 0 {
		return nil, ErrKeysNotFound
	}

	b.logger.Infof("Got keys from SMembers: %v\n", keys)
	data, err := b.redisStore.Client.MGet(keys...).Result()
	if err != nil {
		return nil, ErrRequestNotFound
	}

	reqs := make([]*proto.AlertReq, 0)
	for _, e := range data {
		var req *proto.AlertReq
		if err := json.Unmarshal([]byte(e.(string)), &req); err != nil {
			continue
		}

		reqs = append(reqs, req)
	}

	return reqs, nil
}

func (b *Binance) registerRequest(req *proto.AlertReq) error {
	var membersKey string
	if req.Req.Action == "gt" {
		membersKey = GT
	} else {
		membersKey = LT
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var res *redis.BoolCmd
	_, err = b.redisStore.Client.TxPipelined(func(pipeliner redis.Pipeliner) error {
		res = pipeliner.SetNX(reqKey(req.Req.Price), data, 0)
		pipeliner.SAdd(membersKey, req.Req.Price)
		return nil
	})

	if err != nil {
		return fmt.Errorf("redis: %s", err)
	}

	if !res.Val() {
		return ErrAlreadyRegistered
	}

	return nil
}

func (b *Binance) deregisterRequest(reqPrice float32) error {
	if err := b.redisStore.Client.Del(reqKey(reqPrice)).Err(); err != nil {
		return err
	}

	return nil
}

func reqKey(id float32) string {
	return fmt.Sprintf("binance_%v", id)
}
