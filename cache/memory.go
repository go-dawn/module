package cache

import (
	"sync"
	"time"

	"github.com/go-dawn/dawn/config"
)

type item struct {
	v   []byte
	exp int64
}

type memory struct {
	items      sync.Map
	gcInterval time.Duration
}

func resolveMemory(c *config.Config) Storage {
	return memory{gcInterval: c.GetDuration("GCInterval", time.Second*10)}
}

func (m memory) Has(key string) (bool, error) {
	i, ok := m.items.Load(key)
	if !ok || i.(item).exp < time.Now().Unix() {
		m.items.Delete(key)
		return false, nil
	}

	return true, nil
}

func (m memory) Get(key string, defaultValue ...[]byte) ([]byte, error) {
	v, ok := m.items.Load(key)
	if ok {
		if i := v.(item); i.exp >= time.Now().Unix() {
			return v.([]byte), nil
		} else {
			m.items.Delete(key)
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return nil, nil
}

func (m memory) Many(keys []string) (values [][]byte, err error) {
	for _, k := range keys {
		if v, ok := m.items.Load(k); ok {
			values = append(values, v.([]byte))
		} else {
			values = append(values, nil)
		}
	}

	return
}

func (memory) Put(key string, value []byte, expiration time.Duration) error {
	panic("implement me")
}

func (memory) Pull(key string, defaultValue ...[]byte) ([]byte, error) {
	panic("implement me")
}

func (memory) Increment(key string, by ...int) (int, error) {
	panic("implement me")
}

func (memory) Decrement(key string, value ...int) (int, error) {
	panic("implement me")
}

func (memory) Forever(key string, value []byte) error {
	panic("implement me")
}

func (memory) Remember(key string, expiration time.Duration, defaultValueFunc func() ([]byte, error)) ([]byte, error) {
	panic("implement me")
}

func (memory) RememberForever(key string, defaultValueFunc func() ([]byte, error)) ([]byte, error) {
	panic("implement me")
}

func (memory) Forget(key string) error {
	panic("implement me")
}

func (memory) Flush() error {
	panic("implement me")
}

func (memory) gc() {
	panic("implement me")
}
