package cache

import (
	"context"
	"time"

	"github.com/go-dawn/dawn/config"
	dawnRedis "github.com/go-dawn/dawn/db/redis"
	"github.com/go-redis/redis/v8"
)

var cb = context.Background()

// Cmdable is a wrapper of go-redis Cmdable
type Cmdable interface {
	redis.Cmdable
}

type redisStorage struct {
	db Cmdable
}

func newRedis(c *config.Config) redisStorage {
	return redisStorage{db: dawnRedis.Conn(c.GetString("connection"))}
}

func (s redisStorage) Has(key string) (bool, error) {
	i, err := s.db.Exists(cb, key).Result()
	if err != nil {
		return false, err
	}

	return i != 0, nil
}

func (s redisStorage) Get(key string) (b []byte, err error) {
	b, err = s.db.Get(cb, key).Bytes()
	if err == redis.Nil {
		err = nil
	}

	return
}

func (s redisStorage) GetWithDefault(key string, defaultValue []byte) (b []byte, err error) {
	b, err = s.db.Get(cb, key).Bytes()
	if err == redis.Nil {
		err = nil
		b = defaultValue
	}

	return
}

func (s redisStorage) Many(keys []string) (b [][]byte, err error) {
	var values []interface{}
	if values, err = s.db.MGet(cb, keys...).Result(); err != nil {
		return
	}

	for _, v := range values {
		if v != nil {
			b = append(b, []byte(v.(string)))
		} else {
			b = append(b, nil)
		}
	}

	return
}

func (s redisStorage) Set(key string, value []byte, ttl time.Duration) error {
	return s.db.Set(cb, key, value, ttl).Err()
}

func (s redisStorage) Pull(key string) (b []byte, err error) {
	if b, err = s.db.Get(cb, key).Bytes(); err == nil {
		_, err = s.db.Del(cb, key).Result()
		return
	}

	if err == redis.Nil {
		err = nil
	}

	return
}

func (s redisStorage) PullWithDefault(key string, defaultValue []byte) (b []byte, err error) {
	if b, err = s.db.Get(cb, key).Bytes(); err == nil {
		_, err = s.db.Del(cb, key).Result()
		return
	}

	if err == redis.Nil {
		err = nil
		b = defaultValue
	}

	return
}

func (s redisStorage) Forever(key string, value []byte) error {
	return s.db.Set(cb, key, value, 0).Err()
}

func (s redisStorage) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) (b []byte, err error) {
	if b, err = s.db.Get(cb, key).Bytes(); err == nil {
		return
	}

	if err == redis.Nil {
		if b, err = valueFunc(); err == nil {
			_, err = s.db.Set(cb, key, b, ttl).Result()
		}
	}

	return
}

func (s redisStorage) RememberForever(key string, valueFunc func() ([]byte, error)) (b []byte, err error) {
	if b, err = s.db.Get(cb, key).Bytes(); err == nil {
		return
	}

	if err == redis.Nil {
		if b, err = valueFunc(); err == nil {
			_, err = s.db.Set(cb, key, b, 0).Result()
		}
	}

	return
}

func (s redisStorage) Delete(key string) error {
	return s.db.Del(cb, key).Err()
}

func (s redisStorage) Reset() error {
	return s.db.FlushDB(cb).Err()
}

func (s redisStorage) Close() error {
	return nil
}

func (s redisStorage) gc() {
	return
}
