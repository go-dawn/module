package cache

import (
	"testing"
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Cache_Memory_New(t *testing.T) {
	at := assert.New(t)

	s := newMemory(config.New())

	at.Equal(time.Second*10, s.gcInterval)
	at.NotNil(s.done)
}

func Test_Cache_Memory_Has(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b1, err := s.Has("k1")
	at.Nil(err)
	at.False(b1)

	b2, err := s.Has("k2")
	at.Nil(err)
	at.True(b2)
}

func Test_Cache_Memory_Get(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b1, err := s.Get("k1")
	at.Nil(err)
	at.Nil(b1)

	b2, err := s.Get("k2")
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_GetWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b1, err := s.GetWithDefault("k1", []byte("v11"))
	at.Nil(err)
	at.Equal([]byte("v11"), b1)

	b2, err := s.GetWithDefault("k2", []byte("v22"))
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_Many(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	v, err := s.Many([]string{"k1", "k2"})
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
	s := getMemStorage()

	err := s.Set("k1", []byte("v11"), time.Second)
	at.Nil(err)

	v, ok := s.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(memEntry).data))
	at.InDelta(time.Now().Unix(), v.(memEntry).expiry, 1)
}

func Test_Cache_Memory_Pull(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b2, err := s.Pull("k2")
	at.Nil(err)
	at.Equal("v2", string(b2))

	_, ok := s.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_PullWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b1, err := s.PullWithDefault("k1", []byte("v11"))
	at.Nil(err)
	at.Equal("v11", string(b1))

	_, ok := s.db.Load("k1")
	at.False(ok)

	b2, err := s.PullWithDefault("k2", []byte("v22"))
	at.Nil(err)
	at.Equal("v2", string(b2))

	_, ok = s.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_Forever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	err := s.Forever("k1", []byte("v11"))
	at.Nil(err)

	v, ok := s.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(memEntry).data))
	at.Equal(int64(0), v.(memEntry).expiry)
}

func Test_Cache_Memory_Remember(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b1, err := s.Remember("k1", time.Second, func() (bytes []byte, err error) {
		return []byte("v11"), nil
	})
	at.Nil(err)
	at.Equal("v11", string(b1))

	v, ok := s.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(memEntry).data))
	at.InDelta(time.Now().Unix(), v.(memEntry).expiry, 1)

	b2, err := s.Remember("k2", time.Second, func() (bytes []byte, err error) {
		return []byte("v22"), nil
	})
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_RememberForever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	b1, err := s.RememberForever("k1", func() (bytes []byte, err error) {
		return []byte("v11"), nil
	})
	at.Nil(err)
	at.Equal("v11", string(b1))

	v, ok := s.db.Load("k1")
	at.True(ok)
	at.Equal("v11", string(v.(memEntry).data))
	at.Equal(int64(0), v.(memEntry).expiry)

	b2, err := s.RememberForever("k2", func() (bytes []byte, err error) {
		return []byte("v22"), nil
	})
	at.Nil(err)
	at.Equal("v2", string(b2))
}

func Test_Cache_Memory_Delete(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	err := s.Delete("k2")
	at.Nil(err)

	_, ok := s.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_Reset(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	err := s.Reset()
	at.Nil(err)

	_, ok := s.db.Load("k1")
	at.False(ok)

	_, ok = s.db.Load("k2")
	at.False(ok)
}

func Test_Cache_Memory_Close(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	s := getMemStorage()

	at.Nil(s.Close())

	at.Panics(func() {
		close(s.done)
	})
}

func Test_Cache_Memory_GC(t *testing.T) {
	t.Parallel()

	s := getMemStorage()

	go s.gc()

	assert.Eventually(t, func() bool {
		_, b1 := s.db.Load("k1")
		_, b2 := s.db.Load("k2")
		return !b1 && b2
	}, time.Second, time.Millisecond*10)

	close(s.done)
}

// go test -bench=Benchmark_Cache_Memory_Get -benchmem -count=4
func Benchmark_Cache_Memory_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := getMemStorage()

		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			_, _ = s.Get("k2")
		}
	})
}

func getMemStorage() memStorage {
	s := memStorage{
		gcInterval: time.Millisecond * 10,
		done:       make(chan struct{}),
	}

	s.db.Store("k1", memEntry{[]byte("v1"), time.Now().Add(-time.Minute).Unix()})
	s.db.Store("k2", memEntry{[]byte("v2"), time.Now().Add(time.Minute).Unix()})

	return s
}
