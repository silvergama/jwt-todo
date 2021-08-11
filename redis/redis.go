package redis

import (
	"fmt"
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
	redisHost := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClientReader := createClient(redisHost, redisPassword)
	redisClientWriter := createClient(redisHost, redisPassword)
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

func createClient(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          0, // use default DB
		MaxRetries:  2,
		DialTimeout: time.Duration(2) * time.Second,
		ReadTimeout: time.Duration(1) * time.Second,
		PoolTimeout: time.Duration(15) * time.Second,
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
