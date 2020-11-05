package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cache_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:cache", New().String())
}

func Test_Cache_Module_Init(t *testing.T) {
	m := module{}

	m.Init()()

	// more assertions
}

func Test_Module_Boot(t *testing.T) {
	m := &cacheModule{}

	m.Boot()

	// more assertions
}
