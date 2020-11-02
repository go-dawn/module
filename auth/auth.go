package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/fiberx"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

type authModule struct {
	dawn.Module
}

// New returns the module
func New() dawn.Moduler {
	return &authModule{}
}

func (m *authModule) String() string {
	return "dawn:auth"
}

func (m *authModule) Init() dawn.Cleanup {
	// you can implement me or remove me

	// Read config and init module

	return func() {
		// Put cleanup stuff here if any
	}
}

func (m *authModule) Boot() {
	// you can implement me or remove me
}

func (m *authModule) RegisterRoutes(router fiber.Router) {
	router.Post("/login", m.login)

	router.Use(m.jwt())

	router.Get("/hello", m.hello)
}

func (m *authModule) jwt() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return fiberx.Resp(c, fiber.StatusBadRequest, fiberx.Response{
					Message: "Missing or malformed JWT",
				})
			}

			return fiberx.Resp(c, fiber.StatusUnauthorized, fiberx.Response{
				Message: "Invalid or expired JWT",
			})
		},
	})
}

func (m *authModule) login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return fiber.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = "1"
	claims["exp"] = time.Now().Add(time.Second).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return fiberx.Data(c, fiber.Map{"token": t})
}

func (m *authModule) hello(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return fiberx.Message(c, "Welcome "+name)
}
