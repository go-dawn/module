package cache

import (
	"time"

	"github.com/go-dawn/dawn/config"
	dawnRedis "github.com/go-dawn/dawn/db/redis"
	"github.com/go-redis/redis/v8"
)

// Cmdable is an alias for go-redis Cmdable
type Cmdable = redis.Cmdable

type redisStorage struct {
	db Cmdable
}

func newRedis(c *config.Config) redisStorage {
	return redisStorage{db: dawnRedis.Conn(c.GetString("connection"))}
}

func (s redisStorage) Has(key string) (bool, error) {
	panic("implement me")
}

func (s redisStorage) Get(key string) ([]byte, error) {
	panic("implement me")
}

func (s redisStorage) GetWithDefault(key string, defaultValue []byte) ([]byte, error) {
	panic("implement me")
}

func (s redisStorage) Many(keys []string) ([][]byte, error) {
	panic("implement me")
}

func (s redisStorage) Set(key string, value []byte, ttl time.Duration) error {
	panic("implement me")
}

func (s redisStorage) Pull(key string) ([]byte, error) {
	panic("implement me")
}

func (s redisStorage) PullWithDefault(key string, defaultValue []byte) ([]byte, error) {
	panic("implement me")
}

func (s redisStorage) Forever(key string, value []byte) error {
	panic("implement me")
}

func (s redisStorage) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) ([]byte, error) {
	panic("implement me")
}

func (s redisStorage) RememberForever(key string, valueFunc func() ([]byte, error)) ([]byte, error) {
	panic("implement me")
}

func (s redisStorage) Delete(key string) error {
	panic("implement me")
}

func (s redisStorage) Reset() error {
	panic("implement me")
}

func (s redisStorage) Close() error {
	panic("implement me")
}

func (s redisStorage) gc() {
	panic("implement me")
}
