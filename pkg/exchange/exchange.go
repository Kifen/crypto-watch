package exchange

import (
	"github.com/Kifen/crypto-watch/pkg/util"
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

type Exchange struct {
	Store  *RedisStore
	Logger *logging.Logger
	Srv    *Server
}

func NewExchange(redisUrl, password string) (*Exchange, error) {
	s, err := NewRedisStore(redisUrl, password)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		Store:  s,
		Logger: util.Logger("Exchange"),
		Srv:    NewServer(),
	}, nil
}
