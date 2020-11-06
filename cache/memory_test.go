package cache

import (
	"testing"
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Cache_Memory_New(t *testing.T) {
	at := assert.New(t)

	m := newMemory(config.New())

	at.Equal(time.Second*10, m.gcInterval)
	at.NotNil(m.done)
}

func Test_Cache_Memory_Has(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b1, err := m.Has("k1")
	at.Nil(err)
	at.False(b1)

	b2, err := m.Has("k2")
	at.Nil(err)
	at.True(b2)
}

func Test_Cache_Memory_Get(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b1, err := m.Get("k1")
	at.Nil(err)
	at.Nil(b1)

	b2, err := m.Get("k2")
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_GetWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b1, err := m.GetWithDefault("k1", []byte("v11"))
	at.Nil(err)
	at.Equal([]byte("v11"), b1)

	b2, err := m.GetWithDefault("k2", []byte("v22"))
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_Many(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	v, err := m.Many([]string{"k1", "k2"})
	at.Nil(err)
	at.Len(v, 2)

	if v[0] != nil {
		at.Equal("v2", string(v[0]))
	} else {
		at.Equal("v2", string(v[1]))
	}
}

func Test_Cache_Memory_Set(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	err := m.Set("k1", []byte("v11"), time.Second)
	at.Nil(err)

	v, ok := m.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(entry).data))
	at.InDelta(time.Now().Unix(), v.(entry).expiry, 1)
}

func Test_Cache_Memory_Pull(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b2, err := m.Pull("k2")
	at.Nil(err)
	at.Equal("v2", string(b2))

	_, ok := m.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_PullWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b1, err := m.PullWithDefault("k1", []byte("v11"))
	at.Nil(err)
	at.Equal("v11", string(b1))

	_, ok := m.db.Load("k1")
	at.False(ok)

	b2, err := m.PullWithDefault("k2", []byte("v22"))
	at.Nil(err)
	at.Equal("v2", string(b2))

	_, ok = m.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_Forever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	err := m.Forever("k1", []byte("v11"))
	at.Nil(err)

	v, ok := m.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(entry).data))
	at.Equal(int64(0), v.(entry).expiry)
}

func Test_Cache_Memory_Remember(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b1, err := m.Remember("k1", time.Second, func() (bytes []byte, err error) {
		return []byte("v11"), nil
	})
	at.Nil(err)
	at.Equal("v11", string(b1))

	v, ok := m.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(entry).data))
	at.InDelta(time.Now().Unix(), v.(entry).expiry, 1)

	b2, err := m.Remember("k2", time.Second, func() (bytes []byte, err error) {
		return []byte("v22"), nil
	})
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_RememberForever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	b1, err := m.RememberForever("k1", func() (bytes []byte, err error) {
		return []byte("v11"), nil
	})
	at.Nil(err)
	at.Equal("v11", string(b1))

	v, ok := m.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(entry).data))
	at.Equal(int64(0), v.(entry).expiry)

	b2, err := m.RememberForever("k2", func() (bytes []byte, err error) {
		return []byte("v22"), nil
	})
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_Delete(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	err := m.Delete("k2")
	at.Nil(err)

	_, ok := m.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_Reset(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	err := m.Reset()
	at.Nil(err)

	_, ok := m.db.Load("k1")
	at.False(ok)

	_, ok = m.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_Close(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	m := getMemory()

	err := m.Close()
	at.Nil(err)

	at.Panics(func() {
		select {
		case m.done <- struct{}{}:
		case <-time.NewTimer(time.Second).C:
			at.Fail("should panic")
		}
	})
}

func Test_Cache_Memory_GC(t *testing.T) {
	t.Parallel()

	m := getMemory()

	go m.gc()

	time.Sleep(time.Millisecond * 15)

	close(m.done)

	assert.Eventually(t, func() bool {
		_, b1 := m.db.Load("k1")
		_, b2 := m.db.Load("k2")
		return !b1 && b2
	}, time.Second, time.Millisecond*10)
}

// go test -bench=Benchmark_Cache_Memory_Get -benchmem -count=4
func Benchmark_Cache_Memory_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		m := getMemory()

		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			_, _ = m.Get("k2")
		}
	})
}

func getMemory() memory {
	m := memory{
		gcInterval: time.Millisecond * 10,
		done:       make(chan struct{}),
	}

	m.db.Store("k1", entry{[]byte("v1"), time.Now().Add(-time.Minute).Unix()})
	m.db.Store("k2", entry{[]byte("v2"), time.Now().Add(time.Minute).Unix()})

	return m
}
