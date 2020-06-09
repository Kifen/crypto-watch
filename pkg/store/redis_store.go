package store

import (
	"log"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr, password string) (*RedisStore, error) {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}

	//opt.Password = password
	client := redis.NewClient(opt)
	_, err = client.Ping().Result()
	if err != nil {
		log.Printf("Error creating redis Client.")
		return nil, err
	}

	log.Println("Redis setup complete...")
	return &RedisStore{redis.NewClient(opt)}, nil
}

func (r *RedisStore) RegisterExchange(exchange string) error {
	err := r.Client.SAdd("exchanges", exchange).Err()
	if err != nil {
		return err
	}

	return nil
}
