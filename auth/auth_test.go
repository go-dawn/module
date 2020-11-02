package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:auth", New().String())
}

func Test_Module_Init(t *testing.T) {
	m := &authModule{}

	m.Init()()

	// more assertions
}

func Test_Module_Boot(t *testing.T) {
	m := &authModule{}

	m.Boot()

	// more assertions
}
