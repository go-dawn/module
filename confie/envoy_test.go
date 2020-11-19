package confie

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/go-dawn/module/confie/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockErr = errors.New("fake error")

func Test_Confie_Envoy_Call(t *testing.T) {
	at := assert.New(t)

	t.Run("nil", func(t *testing.T) {
		at.Nil(Call())
	})

	t.Run("non-exist", func(t *testing.T) {
		m.envoys = make(map[string]*Envoy)

		at.Nil(Call("non-exist"))
	})
}

func Test_Confie_Envoy_Make(t *testing.T) {
	at := assert.New(t)

	e, senderBuf, storage := mockEnvoy()

	t.Run("failed to store code", func(t *testing.T) {
		storage.On("Set", "key", mock.Anything, time.Minute).
			Once().Return(mockErr)

		err := e.Make("address", "key")

		at.Equal(mockErr, err)
	})

	t.Run("success", func(t *testing.T) {
		storage.On("Set", "key", mock.Anything, time.Minute).
			Once().Return(nil)

		err := e.Make("address", "key")

		at.Nil(err)
		at.Regexp(`Send [\d]{6} to address`, senderBuf.String())
	})
}

func Test_Confie_Envoy_Verify(t *testing.T) {
	at := assert.New(t)

	e, _, storage := mockEnvoy()

	t.Run("failed to get code", func(t *testing.T) {
		storage.On("Get", "key").
			Once().Return(nil, mockErr)

		err := e.Verify("key", "123456")

		at.Equal(mockErr, err)
	})

	t.Run("not match", func(t *testing.T) {
		storage.On("Get", "key").
			Once().Return([]byte("123456"), nil)

		err := e.Verify("key", "000000")

		at.Equal(ErrNotMatched, err)
	})

	t.Run("success", func(t *testing.T) {
		storage.On("Get", "key").
			Once().Return([]byte("123456"), nil).
			On("Delete", "key").
			Once().Return(nil)

		err := e.Verify("key", "123456")
		at.Nil(err)
	})
}

func mockEnvoy() (*Envoy, *bytes.Buffer, *mocks.Storage) {
	buf := &bytes.Buffer{}
	mockStorage := new(mocks.Storage)
	m := &Module{codeLen: 6, ttl: time.Minute, Config: &Config{Storage: mockStorage}}
	return &Envoy{m: m, Sender: &localSender{out: buf}}, buf, mockStorage
}
