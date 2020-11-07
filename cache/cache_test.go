package cache

import (
	"testing"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Cache_Module_Name(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "dawn:cache", New().String())
}

func Test_Cache_Module_Init(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("default", func(t *testing.T) {
		m := &Module{}

		m.Init()()

		at.Equal(fallback, m.fallback)
		at.Len(m.storage, 1)
	})

	t.Run("invalid driver", func(t *testing.T) {
		config.Load("./", "invalid")
		m := &Module{}

		at.Panics(func() {
			m.Init()
		})
	})
}

func Test_Cache_Module_Build(t *testing.T) {
	at := assert.New(t)

	c := config.New()

	c.Set("driver", "redis")
	at.NotNil(build("redis", c))

	c.Set("driver", "gorm")
	at.NotNil(build("gorm", c))
}

func Test_Cache_Module_Boot(t *testing.T) {
	t.Parallel()

	m := &Module{
		storage: map[string]Cacher{
			fallback: newMemory(config.New()),
		},
	}

	m.Boot()
}

func Test_Cache_Module_Store(t *testing.T) {
	m = &Module{
		fallback: fallback,
		storage: map[string]Cacher{
			fallback: newMemory(config.New()),
		},
	}

	assert.NotNil(t, Storage(fallback))
}
