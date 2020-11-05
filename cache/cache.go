package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-dawn/dawn/config"

	"github.com/go-dawn/dawn"
)

var (
	m        = module{storage: make(map[string]Storage)}
	fallback = "memory"
)

type module struct {
	dawn.Module
	storage  map[string]Storage
	fallback string
}

// New returns the module
func New() dawn.Moduler {
	return m
}

func (m module) String() string {
	return "dawn:cache"
}

func (m module) Init() dawn.Cleanup {
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

	return nil
}

func build(name string, c *config.Config) Storage {
	driver := c.GetString("driver", "memory")

	switch strings.ToLower(driver) {
	case "memory":
		return resolveMemory(c)
	case "redis":
		return resolveRedis(c)
	case "gorm":
		return resolveGorm(c)
	default:
		panic(fmt.Sprintf("dawn:cache unknown driver %s of %s", driver, name))
	}
}

// Storage interface defines cache behaviors.
type Storage interface {
	// Has determines if an item exists in the cache.
	Has(key string) (bool, error)

	// Get retrieves an item from the cache by key.
	// Or use default value if value doesn't exist.
	Get(key string, defaultValue ...[]byte) ([]byte, error)

	// Many retrieves multiple items from the cache by key.
	// Items not found in the cache will have a empty string.
	Many(keys []string) ([][]byte, error)

	// Put stores an item in the cache for a given number of expiration.
	Put(key string, value []byte, expiration time.Duration) error

	// Pull retrieves an item from the cache and delete it.
	// Or use default value if the item doesn't exist.
	Pull(key string, defaultValue ...[]byte) ([]byte, error)

	// Increment increase an item in the cache, default by 1.
	Increment(key string, by ...int) (int, error)

	// Decrement decrease an item in the cache, default by 1.
	Decrement(key string, value ...int) (int, error)

	// Forever stores an item in the cache indefinitely.
	Forever(key string, value []byte) error

	// Remember stores an item from the closure in the cache
	// for a given number of expiration.
	Remember(key string, expiration time.Duration, defaultValueFunc func() ([]byte, error)) ([]byte, error)

	// RememberForever stores an item from the closure in the
	// cache forever.
	RememberForever(key string, defaultValueFunc func() ([]byte, error)) ([]byte, error)

	// Forget removes an item from the cache.
	Forget(key string) error

	// Flush removes all items from the cache.
	Flush() error

	gc()
}

// Store gets cache storage by specific name
// or default.
func Store(name ...string) Storage {
	n := m.fallback

	if len(name) > 0 {
		n = name[0]
	}

	return m.storage[n]
}
