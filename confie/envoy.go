package confie

import (
	"errors"
	"time"

	"github.com/go-dawn/pkg/rand"
)

// ErrNotMatched occurs when code is not matched
var ErrNotMatched = errors.New("confie: code not matched")

// Envoy can generate code and send to a specific address
type Envoy struct {
	// Sender is embedded for sending code
	Sender

	m *Module
}

// Call returns a named envoy
func Call(name ...string) *Envoy {
	if m.envoys == nil {
		return nil
	}

	n := fallback
	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	return m.envoys[n]
}

// Make generates code and sends it
func (e *Envoy) Make(address, key string, ttl time.Duration) (err error) {
	c := rand.NumBytes(e.m.codeLen)

	if err = e.m.Set(key, c, ttl); err != nil {
		return
	}

	err = e.Send(address, string(c))

	return
}

func (e *Envoy) Verify(key, code string) (bool, error) {
	b, err := e.m.Get(key)
	if err != nil {
		return false, err
	}
	if code != string(b) {
		return false, ErrNotMatched
	}

	_ = e.m.Delete(key)

	return true, nil
}
