package cache

import (
	"sync"
	"time"

	"github.com/go-dawn/dawn/config"
)

type memEntry struct {
	data   []byte
	expiry int64
}

type memStorage struct {
	db         sync.Map
	gcInterval time.Duration
	done       chan struct{}
}

func newMemory(c *config.Config) memStorage {
	return memStorage{
		gcInterval: c.GetDuration("GCInterval", time.Second*10),
		done:       make(chan struct{}),
	}
}

func (s memStorage) Has(key string) (bool, error) {
	if v, ok := s.db.Load(key); ok {
		if i := v.(memEntry); i.expiry >= time.Now().Unix() {
			return true, nil
		}
		// Delete expired entry
		s.db.Delete(key)
	}
	return false, nil
}

func (s memStorage) Get(key string) ([]byte, error) {
	return s.value(key), nil
}

func (s memStorage) GetWithDefault(key string, defaultValue []byte) ([]byte, error) {
	v := s.value(key)

	if v == nil {
		v = defaultValue
	}

	return v, nil
}

func (s memStorage) Many(keys []string) (values [][]byte, err error) {
	for _, key := range keys {
		values = append(values, s.value(key))
	}

	return
}

func (s memStorage) Set(key string, value []byte, ttl time.Duration) error {
	s.db.Store(key, memEntry{data: value, expiry: time.Now().Add(ttl).Unix()})
	return nil
}

func (s memStorage) Pull(key string) ([]byte, error) {
	v := s.value(key)
	if v != nil {
		s.db.Delete(key)
	}
	return v, nil
}

func (s memStorage) PullWithDefault(key string, defaultValue []byte) ([]byte, error) {
	v := s.value(key)
	if v != nil {
		s.db.Delete(key)
	} else {
		v = defaultValue
	}
	return v, nil
}

func (s memStorage) Forever(key string, value []byte) error {
	s.db.Store(key, memEntry{data: value, expiry: 0})
	return nil
}

func (s memStorage) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) (v []byte, err error) {
	if v = s.value(key); v == nil {
		if v, err = valueFunc(); err == nil {
			s.db.Store(key, memEntry{data: v, expiry: time.Now().Add(ttl).Unix()})
		}
	}
	return
}

func (s memStorage) RememberForever(key string, valueFunc func() ([]byte, error)) (v []byte, err error) {
	if v = s.value(key); v == nil {
		if v, err = valueFunc(); err == nil {
			s.db.Store(key, memEntry{data: v, expiry: 0})
		}
	}
	return
}

func (s memStorage) Delete(key string) error {
	s.db.Delete(key)
	return nil
}

func (s memStorage) Reset() error {
	s.db.Range(func(key, _ interface{}) bool {
		s.db.Delete(key)
		return true
	})
	return nil
}

func (s memStorage) Close() error {
	close(s.done)
	return nil
}

func (s memStorage) gc() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case t := <-ticker.C:
			now := t.Unix()
			s.db.Range(func(key, value interface{}) bool {
				if e := value.(memEntry); e.expiry != 0 && e.expiry < now {
					s.db.Delete(key)
				}
				return true
			})
		}
	}
}

func (s memStorage) value(key string) []byte {
	if v, ok := s.db.Load(key); ok {
		if e := v.(memEntry); e.expiry == 0 || e.expiry >= time.Now().Unix() {
			return e.data
		}
		// Delete expired entry
		s.db.Delete(key)
	}
	return nil
}
