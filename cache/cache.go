package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-dawn/dawn/config"

	"github.com/go-dawn/dawn"
)

var (
	m        Module
	fallback = "memory"
)

type Module struct {
	dawn.Module
	storage  map[string]Storage
	fallback string
}

// New returns the Module
func New() Module {
	return m
}

func (m Module) String() string {
	return "dawn:cache"
}

func (m Module) Init() dawn.Cleanup {
	m.storage = make(map[string]Storage)

	// extract cache config
	c := config.Sub("cache")

	m.fallback = c.GetString("default", fallback)

	storeConfig := c.GetStringMap("storage")

	if len(storeConfig) == 0 {
		m.fallback = fallback
		m.storage[m.fallback] = build(m.fallback, config.New())
	}

	// build each storage in config
	for name := range storeConfig {
		cfg := c.Sub("storage." + name)
		m.storage[name] = build(name, cfg)
	}

	return m.cleanup
}

func build(name string, c *config.Config) Storage {
	driver := c.GetString("driver", "memory")

	switch strings.ToLower(driver) {
	case "memory":
		return newMemory(c)
	case "redis":
		return newRedis(c)
	case "gorm":
		return newGorm(c)
	default:
		panic(fmt.Sprintf("dawn:cache unknown driver %s of %s", driver, name))
	}
}

func (m Module) Boot() {
	for _, s := range m.storage {
		go s.gc()
	}
}

func (m Module) cleanup() {
	for _, s := range m.storage {
		_ = s.Close()
	}
}

// Storage interface defines cache behaviors.
type Storage interface {
	// Has determines if an entry exists in the cache.
	Has(key string) (bool, error)

	// Get retrieves an entry from the cache for the given key.
	Get(key string) ([]byte, error)

	// GetWithDefault retrieves an entry from the cache for the
	// given key. Returns default value if value is not found.
	GetWithDefault(key string, defaultValue []byte) ([]byte, error)

	// Many retrieves multiple db from the cache by key.
	// Items not found in the cache will have a empty string.
	Many(keys []string) ([][]byte, error)

	// Set stores an entry in the cache for a given number of ttl.
	Set(key string, value []byte, ttl time.Duration) error

	// Pull retrieves an entry from the cache and removes it in the cache.
	Pull(key string) ([]byte, error)

	// PullWithDefault retrieves an entry from the cache and removes it in
	// the cache. Returns default value if value is not found.
	PullWithDefault(key string, defaultValue []byte) ([]byte, error)

	// Forever stores an entry in the cache indefinitely.
	Forever(key string, value []byte) error

	// Remember gets an entry from the cache, or stores an entry from
	// the closure in the cache for a given number of ttl.
	Remember(key string, ttl time.Duration, valueFunc func() ([]byte, error)) ([]byte, error)

	// RememberForever gets an entry from the cache, or stores an entry
	// from the closure in the cache forever.
	RememberForever(key string, valueFunc func() ([]byte, error)) ([]byte, error)

	// Delete removes an entry from the cache.
	Delete(key string) error

	// Reset removes all data from the cache.
	Reset() error

	// Close closes the cache
	Close() error

	gc()
}

// Store gets cache storage by specific name or fallback.
func Store(name ...string) Storage {
	n := m.fallback

	if len(name) > 0 {
		n = name[0]
	}

	return m.storage[n]
}
