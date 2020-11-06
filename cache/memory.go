package cache

import (
	"sync"
	"time"

	"github.com/go-dawn/dawn/config"
)

type entry struct {
	data   []byte
	expiry int64
}

type memory struct {
	db         sync.Map
	gcInterval time.Duration
	done       chan struct{}
}

func newMemory(c *config.Config) memory {
	return memory{
		gcInterval: c.GetDuration("GCInterval", time.Second*10),
		done:       make(chan struct{}),
	}
}

func (m memory) Has(key string) (bool, error) {
	if v, ok := m.db.Load(key); ok {
		if i := v.(entry); i.expiry >= time.Now().Unix() {
			return true, nil
		}
		// Delete expired entry
		m.db.Delete(key)
	}
	return false, nil
}

func (m memory) Get(key string) ([]byte, error) {
	return m.value(key), nil
}

func (m memory) GetWithDefault(key string, defaultValue []byte) ([]byte, error) {
	v := m.value(key)

	if v == nil {
		v = defaultValue
	}

	return v, nil
}

func (m memory) Many(keys []string) (values [][]byte, err error) {
	for _, key := range keys {
		values = append(values, m.value(key))
	}

	return
}

func (m memory) Set(key string, value []byte, ttl time.Duration) error {
	m.db.Store(key, entry{data: value, expiry: time.Now().Add(ttl).Unix()})
	return nil
}

func (m memory) Pull(key string) ([]byte, error) {
	v := m.value(key)
	if v != nil {
		m.db.Delete(key)
	}
	return v, nil
}

func (m memory) PullWithDefault(key string, defaultValue []byte) ([]byte, error) {
	v := m.value(key)
	if v != nil {
		m.db.Delete(key)
	} else {
		v = defaultValue
	}
	return v, nil
}

func (m memory) Forever(key string, value []byte) error {
	m.db.Store(key, entry{data: value, expiry: 0})
	return nil
}

func (m memory) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) (v []byte, err error) {
	if v = m.value(key); v == nil {
		if v, err = valueFunc(); err == nil {
			m.db.Store(key, entry{data: v, expiry: time.Now().Add(ttl).Unix()})
		}
	}
	return
}

func (m memory) RememberForever(key string, valueFunc func() ([]byte, error)) (v []byte, err error) {
	if v = m.value(key); v == nil {
		if v, err = valueFunc(); err == nil {
			m.db.Store(key, entry{data: v, expiry: 0})
		}
	}
	return
}

func (m memory) Delete(key string) error {
	m.db.Delete(key)
	return nil
}

func (m memory) Reset() error {
	m.db.Range(func(key, _ interface{}) bool {
		m.db.Delete(key)
		return true
	})
	return nil
}

func (m memory) Close() error {
	close(m.done)
	return nil
}

func (m memory) gc() {
	ticker := time.NewTicker(m.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.done:
			return
		case t := <-ticker.C:
			now := t.Unix()
			m.db.Range(func(key, value interface{}) bool {
				if e := value.(entry); e.expiry != 0 && e.expiry < now {
					m.db.Delete(key)
				}
				return true
			})
		}
	}
}

func (m memory) value(key string) []byte {
	if v, ok := m.db.Load(key); ok {
		if e := v.(entry); e.expiry == 0 || e.expiry >= time.Now().Unix() {
			return e.data
		}
		// Delete expired entry
		m.db.Delete(key)
	}
	return nil
}
