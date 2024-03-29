package auth

import (
	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/db/sql"
	"github.com/gofiber/fiber/v2"
)

var defaultConfigPath = "config/auth"

type authenticate func(key, code string) (int, error)

type module struct {
	dawn.Module
	*Config
}

// New returns the module
func New(cfg ...*Config) dawn.Moduler {
	var c *Config
	if len(cfg) > 0 {
		c = cfg[0]
	}
	return module{Config: c}
}

func (m module) String() string {
	return "dawn:auth"
}

func (m module) Init() dawn.Cleanup {
	m.setupConfig()

	// Use custom Service
	if m.Service == nil {
		// TODO need code validator implementation
		m.Service = service{repository{sql.Conn()}, nil}
	}

	return nil
}

func (m module) RegisterRoutes(router fiber.Router) {
	g := router.Group("/auth")

	g.Post("/login", m.login)

	g.Use(m.jwt())
}
