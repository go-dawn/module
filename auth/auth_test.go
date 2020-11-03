package auth

import (
	"testing"
	"time"

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
		at.NotNil(m.Repo)
	})
}
