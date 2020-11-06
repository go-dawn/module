package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cache_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:cache", New().String())
}

func Test_Cache_Module_Init(t *testing.T) {
	m := Module{}

	m.Init()

	// more assertions
}

func Test_Cache_Module_Boot(t *testing.T) {
	m := Module{}

	m.Boot()

	// more assertions
}
