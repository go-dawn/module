package confie

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-dawn/dawn/config"
)

type localSender struct {
	out io.Writer
	f   *os.File
}

func newLocalSender(c *config.Config) *localSender {
	l := &localSender{out: os.Stderr}
	if fp := c.GetString("LogFile"); fp != "" {
		if f, err := os.OpenFile(filepath.Clean(fp), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600); err == nil {
			l.f = f
			l.out = f
		}
	}

	return l
}

func (s *localSender) Send(address, code string) error {
	_, err := fmt.Fprintf(s.out, "Send %s to %s\n", code, address)
	return err
}

func (s *localSender) Close() error {
	if s.f != nil {
		return s.f.Close()
	}

	return nil
}
