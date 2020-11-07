package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/module/cache/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var mockErr = errors.New("fake error")

func Test_Cache_Redis_New(t *testing.T) {
	t.Parallel()

	s := newRedis(config.New())
	assert.Nil(t, s.db)
}

func Test_Cache_Redis_Has(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Exists", cb, "k1").
			Once().Return(redis.NewIntResult(0, nil))

		ok, err := s.Has("k1")
		at.Nil(err)
		at.False(ok)
	})

	t.Run("error", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Exists", cb, "k1").
			Once().Return(redis.NewIntResult(0, mockErr))

		_, err := s.Has("k1")
		at.Equal(mockErr, err)
	})
}

func Test_Cache_Redis_Get(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", nil))

		b, err := s.Get("k1")
		at.Nil(err)
		at.Equal("v1", string(b))
	})

	t.Run("non-exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("", redis.Nil))

		b, err := s.Get("k1")
		at.Nil(err)
		at.Nil(b)
	})

	t.Run("error", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", mockErr))

		_, err := s.Get("k1")
		at.Equal(mockErr, err)
	})
}

func Test_Cache_Redis_GetWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", nil))

		b, err := s.GetWithDefault("k1", []byte("v11"))
		at.Nil(err)
		at.Equal("v1", string(b))
	})

	t.Run("non-exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("", redis.Nil))

		b, err := s.GetWithDefault("k1", []byte("v11"))
		at.Nil(err)
		at.Equal("v11", string(b))
	})

	t.Run("error", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", mockErr))

		_, err := s.Get("k1")
		at.Equal(mockErr, err)
	})
}

func Test_Cache_Redis_Many(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("MGet", cb, "k1", "k2").
			Once().Return(redis.NewSliceResult([]interface{}{nil, "v2"}, nil))

		b, err := s.Many([]string{"k1", "k2"})
		at.Nil(err)
		at.Len(b, 2)
		at.Nil(b[0])
		at.Equal("v2", string(b[1]))
	})

	t.Run("error", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("MGet", cb, "k1", "k2").
			Once().Return(redis.NewSliceResult(nil, mockErr))

		_, err := s.Many([]string{"k1", "k2"})
		at.Equal(mockErr, err)
	})
}

func Test_Cache_Redis_Set(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Set", cb, "k1", []byte("v1"), time.Second*10).
			Once().Return(redis.NewStatusResult("OK", nil))

		err := s.Set("k1", []byte("v1"), time.Second*10)
		at.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Set", cb, "k1", []byte("v1"), time.Second*10).
			Once().Return(redis.NewStatusResult("", mockErr))

		err := s.Set("k1", []byte("v1"), time.Second*10)
		at.Equal(mockErr, err)
	})
}

func Test_Cache_Redis_Pull(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", nil)).
			On("Del", cb, "k1").
			Once().Return(redis.NewIntResult(1, nil))

		b, err := s.Pull("k1")
		at.Nil(err)
		at.Equal("v1", string(b))
	})

	t.Run("non-exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("", redis.Nil))

		b, err := s.Pull("k1")
		at.Nil(err)
		at.Nil(b)
	})
}

func Test_Cache_Redis_PullWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", nil)).
			On("Del", cb, "k1").
			Once().Return(redis.NewIntResult(1, nil))

		b, err := s.PullWithDefault("k1", []byte("v11"))
		at.Nil(err)
		at.Equal("v1", string(b))
	})

	t.Run("non-exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("", redis.Nil))

		b, err := s.PullWithDefault("k1", []byte("v11"))
		at.Nil(err)
		at.Equal("v11", string(b))
	})
}

func Test_Cache_Redis_Forever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s, mockDB := getRedisStorage()

	mockDB.On("Set", cb, "k1", []byte("v1"), time.Duration(0)).
		Once().Return(redis.NewStatusResult("OK", nil))

	err := s.Forever("k1", []byte("v1"))
	at.Nil(err)
}

func Test_Cache_Redis_Remember(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", nil))

		b, err := s.Remember("k1", time.Second, func() ([]byte, error) {
			return []byte("v11"), nil
		})
		at.Nil(err)
		at.Equal("v1", string(b))
	})

	t.Run("non-exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("", redis.Nil)).
			On("Set", cb, "k1", []byte("v11"), time.Second).
			Once().Return(redis.NewStatusResult("OK", nil))

		b, err := s.Remember("k1", time.Second, func() ([]byte, error) {
			return []byte("v11"), nil
		})
		at.Nil(err)
		at.Equal("v11", string(b))
	})
}

func Test_Cache_Redis_RememberForever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("v1", nil))

		b, err := s.RememberForever("k1", func() ([]byte, error) {
			return []byte("v11"), nil
		})
		at.Nil(err)
		at.Equal("v1", string(b))
	})

	t.Run("non-exist", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Get", cb, "k1").
			Once().Return(redis.NewStringResult("", redis.Nil)).
			On("Set", cb, "k1", []byte("v11"), time.Duration(0)).
			Once().Return(redis.NewStatusResult("OK", nil))

		b, err := s.RememberForever("k1", func() ([]byte, error) {
			return []byte("v11"), nil
		})
		at.Nil(err)
		at.Equal("v11", string(b))
	})
}

func Test_Cache_Redis_Delete(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s, mockDB := getRedisStorage()

	mockDB.On("Del", cb, "k1").
		Once().Return(redis.NewIntResult(1, nil))

	err := s.Delete("k1")
	at.Nil(err)
}

func Test_Cache_Redis_Reset(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Scan", cb, uint64(0), "*", int64(1000)).
			Once().Return(redis.NewScanCmdResult([]string{"k1", "k2"}, 10, nil)).
			On("Scan", cb, uint64(10), "*", int64(1000)).
			Once().Return(redis.NewScanCmdResult([]string{}, 20, nil)).
			On("Del", cb, "k1", "k2").
			Once().Return(redis.NewIntResult(2, nil))

		at.Nil(s.Reset())
	})

	t.Run("error", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Scan", cb, uint64(0), "*", int64(1000)).
			Once().Return(redis.NewScanCmdResult(nil, 10, mockErr))

		at.Equal(mockErr, s.Reset())
	})
}

func Test_Cache_Redis_Close(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s, _ := getRedisStorage()

	err := s.Close()
	at.Nil(err)
}

func Test_Cache_Redis_GC(t *testing.T) {
	t.Parallel()

	(redisStorage{}).gc()
}

func getRedisStorage() (redisStorage, *mocks.Cmdable) {
	mockDB := new(mocks.Cmdable)
	return redisStorage{db: mockDB}, mockDB
}

func TestRedis(t *testing.T) {
	c := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	t.Log(c.Scan(cb, 0, "dawn_*", 1).Result())
}
