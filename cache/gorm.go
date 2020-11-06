package cache

import (
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/dawn/db/sql"
	"gorm.io/gorm"
)

type gormStorage struct {
	db *gorm.DB
}

func newGorm(c *config.Config) gormStorage {
	return gormStorage{db: sql.Conn(c.GetString("connection"))}
}

func (s gormStorage) Has(key string) (bool, error) {
	panic("implement me")
}

func (s gormStorage) Get(key string) ([]byte, error) {
	panic("implement me")
}

func (s gormStorage) GetWithDefault(key string, defaultValue []byte) ([]byte, error) {
	panic("implement me")
}

func (s gormStorage) Many(keys []string) ([][]byte, error) {
	panic("implement me")
}

func (s gormStorage) Set(key string, value []byte, ttl time.Duration) error {
	panic("implement me")
}

func (s gormStorage) Pull(key string) ([]byte, error) {
	panic("implement me")
}

func (s gormStorage) PullWithDefault(key string, defaultValue []byte) ([]byte, error) {
	panic("implement me")
}

func (s gormStorage) Forever(key string, value []byte) error {
	panic("implement me")
}

func (s gormStorage) Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) ([]byte, error) {
	panic("implement me")
}

func (s gormStorage) RememberForever(key string, valueFunc func() ([]byte, error)) ([]byte, error) {
	panic("implement me")
}

func (s gormStorage) Delete(key string) error {
	panic("implement me")
}

func (s gormStorage) Reset() error {
	panic("implement me")
}

func (s gormStorage) Close() error {
	panic("implement me")
}

func (s gormStorage) gc() {
	panic("implement me")
}
