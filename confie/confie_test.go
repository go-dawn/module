package confie

import (
	"testing"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Confie_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:confie", New(&Config{}).String())
}

func Test_Confie_Module_Init(t *testing.T) {
	m := &Module{}

	at := assert.New(t)

	t.Run("default", func(t *testing.T) {
		cleanup := m.Init()

		at.Equal(6, m.codeLen)
		at.Len(m.envoys, 1)

		cleanup()
	})

	t.Run("invalid driver", func(t *testing.T) {
		config.Load("./", "invalid")
		m := &Module{}

		at.Panics(func() {
			m.Init()
		})
	})
}
