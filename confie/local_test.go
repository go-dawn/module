package confie

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Confie_Local_New(t *testing.T) {
	at := assert.New(t)

	f, err := ioutil.TempFile("", "")
	at.Nil(err)
	defer func() { _ = f.Close() }()

	c := config.New()
	c.Set("LogFile", f.Name())

	l := newLocalSender(c)

	at.Nil(l.Close())
}

func Test_Confie_Local_Send(t *testing.T) {
	at := assert.New(t)

	var b bytes.Buffer
	l := &localSender{out: &b}

	err := l.Send("address", "123456")
	at.Nil(err)

	at.Equal("Send 123456 to address\n", b.String())
}
