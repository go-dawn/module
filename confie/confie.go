package confie

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/module/cache"
)

var m = &Module{}
var fallback = "local"

type Module struct {
	dawn.Module
	*Config
	envoys  map[string]*Envoy
	codeLen int
	ttl     time.Duration
}

// New returns the Module
func New(c ...*Config) *Module {
	if len(c) > 0 {
		m.Config = c[0]
	}
	return m
}

func (m *Module) String() string {
	return "dawn:confie"
}

func (m *Module) Init() dawn.Cleanup {
	m.envoys = make(map[string]*Envoy)

	if m.Config == nil {
		m.Config = &Config{}
	}

	// extract confie config
	c := config.Sub("confie")

	m.codeLen = c.GetInt("codeLength", 6)
	m.ttl = c.GetDuration("ttl", time.Minute*5)

	if m.Storage == nil {
		m.Storage = cache.Storage()
	}

	m.setupEnvoys(c)

	return m.cleanup
}

func (m *Module) setupEnvoys(c *config.Config) {
	fallback := c.GetString("default", fallback)

	drivers := c.GetStringMap("envoys")

	if len(drivers) == 0 {
		m.envoys[fallback] = m.build(fallback, config.New())
	}

	// build each storage in config
	for name := range drivers {
		cfg := c.Sub("envoys." + name)
		m.envoys[name] = m.build(name, cfg)
	}
}

func (m *Module) build(name string, c *config.Config) *Envoy {
	var sender Sender
	switch strings.ToLower(name) {
	case "local":
		sender = newLocalSender(c)
	default:
		panic(fmt.Sprintf("dawn:confie unknown driver %s", name))
	}

	return &Envoy{m: m, Sender: sender}
}

func (m *Module) cleanup() {
	for _, e := range m.envoys {
		e.Sender.close()
	}
}

// Sender defines an interface to send code with expiration
type Sender interface {
	// Send delivers code to the address
	Send(address, code string) error

	close()
}

// Storage defines behaviors to manage sent code
type Storage interface {
	// Get retrieves code from the storage for the given key.
	Get(key string) ([]byte, error)
	// Set stores code in storage with given ttl
	Set(key string, value []byte, ttl time.Duration) error
	// Delete removes code in storage
	Delete(key string) error
}
