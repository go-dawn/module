package auth

import (
	"time"

	"github.com/go-dawn/dawn/config"
)

// Config defines the config for auth module
type Config struct {
	// Repo is a custom repository for auth module
	Repo

	// SigningKey is for generating and validating jwt token
	SigningKey string

	// SigningMethod is used to check token signing method
	// Optional. Default: "HS256"
	// Possible values: "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "RS256", "RS384", "RS512"
	SigningMethod string

	// Expiration is the effective duration of jwt token
	Expiration time.Duration
}

func (m module) setupConfig() {
	if m.Config == nil {
		m.Config = &Config{}

		// Try to read from global Config's auth section
		_ = config.Sub("auth").Unmarshal(m.Config)

		// Try to read from Config/auth
		if m.SigningKey == "" {
			_ = config.New(defaultConfigPath).Unmarshal(m.Config)
		}
	}

	if m.SigningKey == "" {
		panic("auth: signing key is required")
	}

	if m.Expiration == 0 {
		m.Expiration = time.Hour
	}
}
