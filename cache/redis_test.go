package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/go-dawn/module/cache/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var mockErr = errors.New("fake error")

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

func Test_Cache_Redis_Set(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		s, mockDB := getRedisStorage()

		mockDB.On("Set", cb, "k1", "v2").
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

func Test_Cache_Redis_Many(t *testing.T) {
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

func TestRedis(t *testing.T) {
	c := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	b, _ := c.Set(cb, "k1", []byte("k2"), time.Second*10).Result()

	t.Log(b)
}

func getRedisStorage() (redisStorage, *mocks.Cmdable) {
	mockDB := new(mocks.Cmdable)
	return redisStorage{db: mockDB}, mockDB
}
