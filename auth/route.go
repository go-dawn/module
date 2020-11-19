package auth

import (
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/go-dawn/dawn/fiberx"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func (m module) jwt() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(m.SigningKey),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return fiberx.CodeErr(fiber.StatusBadRequest, err)
			}

			return fiberx.CodeErr(fiber.StatusBadRequest, err, "Invalid or expired JWT")
		},
	})
}

type loginForm struct {
	// Username is account username, mobile number or email address
	Username string `json:"username" validate:"required"`
	// Type can be password, mobile or email
	Type string `json:"type" validate:"required,oneof=password mobile email"`
	// Code can be password, sms code or email code
	Code string `json:"code" validate:"required"`
}

func (m module) login(c *fiber.Ctx) (err error) {
	var (
		data loginForm
		id   int
		t    string
	)

	if err = fiberx.ValidateBody(c, &data); err != nil {
		return
	}

	if id, err = m.authFunc(data.Type)(data.Username, data.Code); err != nil {
		return fiberx.CodeErr(fiber.StatusUnauthorized, err, "Failed to authenticate")
	}

	// Generate encoded token and send it as response.
	if t, err = generateToken(m.SigningMethod, m.SigningKey, id, m.Expiration); err != nil {
		return err
	}

	return fiberx.Data(c, t)
}

func (m module) authFunc(tpy string) authenticate {
	switch tpy {
	case "password":
		return m.LoginByPassword
	case "mobile":
		return m.LoginByMobileCode
	case "email":
		return m.LoginByEmailCode
	default:
		return func(username, code string) (i int, err error) {
			return 0, fmt.Errorf("auth: invalid authenticate type %s", tpy)
		}
	}
}

func signingMethod(method string) jwt.SigningMethod {
	switch method {
	case "", "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	case "ES256":
		return jwt.SigningMethodES256
	case "ES384":
		return jwt.SigningMethodES384
	case "ES512":
		return jwt.SigningMethodES512
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS384":
		return jwt.SigningMethodRS384
	case "RS512":
		return jwt.SigningMethodRS512
	default:
		return jwt.SigningMethodHS256
	}
}

func generateToken(method, key string, id int, expiration time.Duration) (t string, err error) {
	// Create token
	token := jwt.New(signingMethod(method))

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(expiration).Unix()

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(key))
}
