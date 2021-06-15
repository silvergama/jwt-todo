package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisService interface {
	Set(key string, value []byte, exp time.Duration) error
	Get(key string) ([]byte, error)
	Del(key string) int64
	FlushAll() error
}

type service struct {
	clientReader *redis.Client
	clientWriter *redis.Client
}

var cacheInstance *service

func GetCacheInstance() RedisService {
	return cacheInstance
}

func Setup() error {
	redisClientReader := createClient("localhost:6379")
	redisClientWriter := createClient("localhost:6379")
	cacheInstance = &service{
		clientReader: redisClientReader,
		clientWriter: redisClientWriter,
	}

	_, err := redisClientReader.Ping().Result()
	if err != nil {
		return err
	}

	_, err = redisClientWriter.Ping().Result()
	return err
}

func createClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func (s *service) Set(key string, value []byte, exp time.Duration) error {
	return s.clientWriter.Set(key, value, exp).Err()
}

func (s *service) Get(key string) ([]byte, error) {
	return s.clientReader.Get(key).Bytes()
}

func (s *service) Del(key string) int64 {
	return s.clientWriter.Del(key).Val()
}

func (s *service) FlushAll() error {
	_, err := s.clientWriter.FlushAll().Result()
	return err
}
