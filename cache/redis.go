package cache

import (
	"context"
	"time"

	"github.com/go-dawn/dawn/config"
	dawnRedis "github.com/go-dawn/dawn/db/redis"
	"github.com/go-redis/redis/v8"
)

var cb = context.Background()

// Cmdable is a subset of go-redis Cmdable
type Cmdable interface {
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
}

type redisStorage struct {
	db     Cmdable
	prefix string
}

func newRedis(c *config.Config) redisStorage {
	return redisStorage{
		db:     dawnRedis.Conn(c.GetString("connection")),
		prefix: c.GetString("prefix", "dawn_cache_"),
	}
}

func (s redisStorage) Has(key string) (bool, error) {
	i, err := s.db.Exists(cb, s.prefixedKey(key)).Result()
	if err != nil {
		return false, err
	}

	return i != 0, nil
}

func (s redisStorage) Get(key string) (b []byte, err error) {
	b, err = s.db.Get(cb, s.prefixedKey(key)).Bytes()
	if err == redis.Nil {
		err = nil
	}

	return
}

func (s redisStorage) GetWithDefault(key string, defaultValue []byte) (b []byte, err error) {
	b, err = s.db.Get(cb, s.prefixedKey(key)).Bytes()
	if err == redis.Nil {
		err = nil
		b = defaultValue
	}

	return
}

func (s redisStorage) Many(keys []string) (b [][]byte, err error) {
	for i := 0; i < len(keys); i++ {
		keys[i] = s.prefixedKey(keys[i])
	}
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
	return s.db.Set(cb, s.prefixedKey(key), value, ttl).Err()
}

func (s redisStorage) Pull(key string) (b []byte, err error) {
	key = s.prefixedKey(key)
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
	key = s.prefixedKey(key)
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
	return s.db.Set(cb, s.prefixedKey(key), value, 0).Err()
}

func (s redisStorage) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) (b []byte, err error) {
	key = s.prefixedKey(key)
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
	key = s.prefixedKey(key)
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
	return s.db.Del(cb, s.prefixedKey(key)).Err()
}

func (s redisStorage) Reset() (err error) {
	var (
		keys    []string
		matched []string
		cursor  uint64
	)

	for {
		if matched, cursor, err = s.db.Scan(cb, cursor, s.prefix+"*", 1000).Result(); err != nil {
			return
		}
		if len(matched) == 0 {
			break
		}
		keys = append(keys, matched...)
	}

	return s.db.Del(cb, keys...).Err()
}

func (s redisStorage) Close() error {
	return nil
}

func (s redisStorage) gc() {
	return
}

func (s redisStorage) prefixedKey(key string) string {
	return s.prefix + key
}
