package auth

import (
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-dawn/dawn/fiberx"
	"github.com/go-dawn/module/auth/mocks"
	"github.com/go-dawn/pkg/deck"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Auth_Route_Login(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	m, mockRepo := routeModule()

	app, e := deck.SetupServer(t)

	app.Post("/login", m.login)

	var (
		username = "kiyon"
		code     = "pass"
		typ      = "password"
	)

	t.Run("bad request", func(t *testing.T) {
		resp := e.POST("/login").Expect()

		resp.Status(fiber.StatusBadRequest)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.On("LoginByPassword", username, code).
			Once().Return(1, nil)

		resp := e.POST("/login").WithJSON(loginForm{
			Username: username,
			Code:     code,
			Type:     typ,
		}).Expect()

		resp.Status(fiber.StatusOK)

		deck.AssertRespDataCheck(resp, func(v *httpexpect.Value) {
			b, err := base64.StdEncoding.DecodeString(v.Raw().(string))
			at.NotNil(err)
			at.Contains(string(b), `{"alg":"HS256","typ":"JWT"}`)
		})
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockRepo.On("LoginByPassword", username, code).
			Once().Return(0, errors.New("invalid"))

		resp := e.POST("/login").WithJSON(loginForm{
			Username: username,
			Code:     code,
			Type:     typ,
		}).Expect()

		resp.Status(fiber.StatusUnauthorized)
		deck.AssertRespMsg(resp, "Failed to authenticate")
	})
}

func Test_Auth_Route_Jwt_Middleware(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	m, _ := routeModule()

	app, e := deck.SetupServer(t)

	app.Get("/", m.jwt(), func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		at.Equal(1, int(claims["id"].(float64)))

		return fiberx.Message(c, "JWT")
	})

	t.Run("missing jwt", func(t *testing.T) {
		resp := e.GET("/").Expect().Status(fiber.StatusBadRequest)

		deck.AssertRespMsg(resp, "Missing or malformed JWT")
	})

	t.Run("invalid jwt", func(t *testing.T) {
		resp := e.GET("/").
			WithHeader(fiber.HeaderAuthorization, "Bearer xxx").
			Expect().
			Status(fiber.StatusBadRequest)

		deck.AssertRespMsg(resp, "Invalid or expired JWT")
	})

	t.Run("success", func(t *testing.T) {
		token, err := generateToken("", m.SigningKey, 1, time.Hour)
		at.Nil(err)

		resp := e.GET("/").
			WithHeader(fiber.HeaderAuthorization, "Bearer "+token).
			Expect().
			Status(fiber.StatusOK)

		deck.AssertRespMsg(resp, "JWT")
	})
}

func Test_Auth_Module_AuthFunc(t *testing.T) {
	at := assert.New(t)

	m, _ := routeModule()

	at.NotNil(m.authFunc("mobile"))
	at.NotNil(m.authFunc("email"))

	fn := m.authFunc("invalid")
	at.NotNil(fn)

	_, err := fn("", "")
	at.NotNil(err)
}

func Test_Auth_SigningMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		sm     jwt.SigningMethod
	}{
		{"", jwt.SigningMethodHS256},
		{"invalid", jwt.SigningMethodHS256},
		{"HS256", jwt.SigningMethodHS256},
		{"HS384", jwt.SigningMethodHS384},
		{"HS512", jwt.SigningMethodHS512},
		{"ES256", jwt.SigningMethodES256},
		{"ES384", jwt.SigningMethodES384},
		{"ES512", jwt.SigningMethodES512},
		{"RS256", jwt.SigningMethodRS256},
		{"RS384", jwt.SigningMethodRS384},
		{"RS512", jwt.SigningMethodRS512},
	}

	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			assert.Equal(t, tc.sm, signingMethod(tc.method))
		})
	}
}

func routeModule() (module, *mocks.Service) {
	mockService := new(mocks.Service)
	return module{Config: &Config{
		Service:    mockService,
		SigningKey: "test",
	}}, mockService
}
