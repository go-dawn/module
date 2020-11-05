package auth

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	defaultConfigPath = "testdata/auth"
}

func Test_Auth_New(t *testing.T) {
	assert.NotNil(t, New(&Config{}))
}

func Test_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:auth", New().String())
}

func Test_Module_Init(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("signing key required", func(t *testing.T) {
		at.Panics(func() {
			(module{}).Init()
		})
	})

	// more assertions
	t.Run("success", func(t *testing.T) {
		m := module{Config: &Config{
			SigningKey: "xx",
		}}

		at.Nil(m.Init())
		at.Equal("xx", m.SigningKey)
		at.Equal(time.Hour, m.Expiration)
		at.NotNil(m.Service)
	})
}

func Test_Module_RegisterRoutes(t *testing.T) {
	t.Parallel()

	m := module{Config: &Config{SigningKey: "xx"}}

	app := fiber.New()

	m.RegisterRoutes(app)

	assertHasRouteGroup(t, app, "/auth")
	assertHasRoute(t, app, fiber.MethodPost, "/auth/login")
}

func assertHasRoute(t *testing.T, app *fiber.App, method string, path string) {
	for _, routes := range app.Stack() {
		for _, r := range routes {
			if r.Method == method && r.Path == path {
				return
			}
		}
	}

	assert.Failf(t, "%s %s not found", method, path)
}

func assertHasRouteGroup(t *testing.T, app *fiber.App, path string) {
	for _, routes := range app.Stack() {
		var found bool
		for _, r := range routes {
			if r.Path == path {
				found = true
				break
			}
		}
		if !found {
			assert.Failf(t, "group %s not found", path)
		}
	}
}
