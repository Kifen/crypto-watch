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
	reqChs       map[string]*proto.AlertReq
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
		reqChs:       make(map[string]*proto.AlertReq),
	}, nil
}

func (b *Binance) subscribe(id string, req *proto.AlertReq) {
	b.rm.Lock()
	if _, ok := b.reqChs[id]; !ok {
		b.reqChs[id] = req
	}
	b.rm.Unlock()
}

func (b *Binance) unsubscribe(id string) {
	delete(b.reqChs, id)
	i, err := strconv.Atoi(id)
	if err != nil {
		b.logger.Fatal(err)
	}

	if err := b.deregisterRequest(int32(i)); err != nil {
		b.logger.Fatalf("failed to deregister request")
	}
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

func (b *Binance) initReqsSubscriptions() error {
	err := b.subRequests(GT)
	if err != nil {
		return err
	}

	err = b.subRequests(LT)
	if err != nil {
		return err
	}

	b.logger.Infof("Added all requests to subscriptions: %v", b.reqChs)
	return nil
}

func (b *Binance) subRequests(key string) error {
	var wg sync.WaitGroup
	reqs, err := b.redisStore.Client.SMembers(key).Result()
	if err != nil {
		return err
	}

	for _, r := range reqs {
		var req *proto.AlertReq
		if err := json.Unmarshal([]byte(r), &req); err != nil {
			b.logger.Fatalf("failed to unmarshall request: %v", err)
		}

		wg.Add(1)
		go b.initSubscribe(req, &wg)
	}
	wg.Wait()

	return nil
}

func (b *Binance) registerRequest(req *proto.AlertReq) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = b.redisStore.Client.Set(reqKey(req.Id), data, 0).Err()
	if err != nil {
		return fmt.Errorf("redis: %s", err)
	}

	return nil
}

func (b *Binance) deregisterRequest(id int32) error {
	if err := b.redisStore.Client.Del(reqKey(id)).Err(); err != nil {
		return err
	}

	return nil
}

func reqKey(id int32) string {
	return fmt.Sprintf("binance_%v", id)
}

func (b *Binance) initSubscribe(r *proto.AlertReq, wg *sync.WaitGroup) {
	defer wg.Done()
	b.subscribe(string(r.Id), r)
}
