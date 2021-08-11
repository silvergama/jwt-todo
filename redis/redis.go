package redis

import (
	"os"
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

	redisClientReader, err := createClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}
	redisClientWriter, err := createClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}
	cacheInstance = &service{
		clientReader: redisClientReader,
		clientWriter: redisClientWriter,
	}

	_, err = redisClientReader.Ping().Result()
	if err != nil {
		return err
	}

	_, err = redisClientWriter.Ping().Result()
	return err
}

func createClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
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
