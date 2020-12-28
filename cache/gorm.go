package cache

import (
	"time"

	"gorm.io/gorm/clause"

	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/dawn/db/sql"
	"gorm.io/gorm"
)

type gormStorage struct {
	db         *gorm.DB
	table      string
	prefix     string
	gcInterval time.Duration
	done       chan struct{}
}

func newGorm(c *config.Config) *gormStorage {
	s := &gormStorage{
		db:         sql.Conn(c.GetString("connection")),
		table:      c.GetString("table", "dawn_cache"),
		prefix:     c.GetString("prefix", "dawn_cache_"),
		gcInterval: c.GetDuration("GCInterval", time.Second*10),
		done:       make(chan struct{}),
	}
	return s.setup()
}

func (s *gormStorage) setup() *gormStorage {
	if s.db != nil {
		s.db = s.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "expiry"}),
		}).
			Table(s.table).
			Session(&gorm.Session{})

		_ = s.db.AutoMigrate(&gormEntry{})
	}

	return s
}

func (s *gormStorage) Has(key string) (ok bool, err error) {
	var e gormEntry
	if err = s.db.Debug().First(&e, "key = ?", s.prefixedKey(key)).Error; err == nil {
		ok = e.valid()
		return
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (s *gormStorage) Get(key string) (b []byte, err error) {
	return s.value(key)
}

func (s *gormStorage) GetWithDefault(key string, defaultValue []byte) (b []byte, err error) {
	if b, err = s.value(key); err == nil && b == nil {
		b = defaultValue
	}

	return
}

func (s *gormStorage) Many(keys []string) (b [][]byte, err error) {
	l := len(keys)
	in := make([]string, l)
	for i := 0; i < l; i++ {
		in[i] = s.prefixedKey(keys[i])
	}

	var entries []gormEntry
	if err = s.db.Find(&entries, "key in ?", in).Error; err == nil {
		b = make([][]byte, l)
		for i := 0; i < l; i++ {
			for _, e := range entries {
				if e.Key == in[i] && e.valid() {
					b[i] = entries[i].Value
					break
				}
			}
		}
	}

	return
}

func (s *gormStorage) Set(key string, value []byte, ttl time.Duration) error {
	return s.set(key, value, time.Now().Add(ttl).Unix())
}

func (s *gormStorage) Pull(key string) (b []byte, err error) {
	if b, err = s.value(key); err == nil && b != nil {
		err = s.Delete(key)
	}

	return
}

func (s *gormStorage) PullWithDefault(key string, defaultValue []byte) (b []byte, err error) {
	if b, err = s.value(key); err == nil {
		if b != nil {
			err = s.Delete(key)
		} else {
			b = defaultValue
		}
	}

	return
}

func (s *gormStorage) Forever(key string, value []byte) error {
	return s.set(key, value, 0)
}

func (s *gormStorage) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) (b []byte, err error) {
	if b, err = s.value(key); err == nil && b == nil {
		if b, err = valueFunc(); err == nil {
			err = s.set(key, b, time.Now().Add(ttl).Unix())
		}
	}

	return
}

func (s *gormStorage) RememberForever(key string, valueFunc func() ([]byte, error)) (b []byte, err error) {
	if b, err = s.value(key); err == nil && b == nil {
		if b, err = valueFunc(); err == nil {
			err = s.set(key, b, 0)
		}
	}

	return
}

func (s *gormStorage) Delete(key string) error {
	return s.db.Delete(&gormEntry{}, "key = ?", s.prefixedKey(key)).Error
}

func (s *gormStorage) Reset() error {
	return s.db.Delete(&gormEntry{}, "key like ?", s.prefix+"%").Error
}

func (s *gormStorage) Close() error {
	close(s.done)
	return nil
}

func (s *gormStorage) gc() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case t := <-ticker.C:
			s.db.Delete(&gormEntry{}, "expiry < ?", t.Unix())
		}
	}
}

func (s *gormStorage) prefixedKey(key string) string {
	return s.prefix + key
}

func (s *gormStorage) value(key string) (b []byte, err error) {
	var e gormEntry
	if err = s.db.First(&e, "key = ?", s.prefixedKey(key)).Error; err == nil {
		if e.valid() {
			b = e.Value
			return
		}
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (s *gormStorage) set(key string, value []byte, expiry int64) error {
	return s.db.Create(&gormEntry{
		Key:    s.prefixedKey(key),
		Value:  value,
		Expiry: expiry,
	}).Error
}

type gormEntry struct {
	Key    string `gorm:"primarykey"`
	Value  []byte
	Expiry int64
}

func (e gormEntry) valid() bool {
	return e.Expiry == 0 || e.Expiry >= time.Now().Unix()
}
