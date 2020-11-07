package cache

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/go-dawn/pkg/deck"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Cache_Gorm_New(t *testing.T) {
	t.Parallel()

	s := newGorm(config.New())
	assert.Nil(t, s.db)
}

func Test_Cache_Gorm_Has(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.Has("k1")
	at.Nil(err)
	at.True(b1)

	b2, err := s.Has("k2")
	at.Nil(err)
	at.False(b2)

	b3, err := s.Has("k3")
	at.Nil(err)
	at.False(b3)
}

func Test_Cache_Gorm_Get(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.Get("k1")
	at.Nil(err)
	at.Equal("v1", string(b1))

	b2, err := s.Get("k2")
	at.Nil(err)
	at.Nil(b2)

	b3, err := s.Get("k3")
	at.Nil(err)
	at.Nil(b3)
}

func Test_Cache_Gorm_GetWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.GetWithDefault("k1", []byte("v11"))
	at.Nil(err)
	at.Equal("v1", string(b1))

	b2, err := s.GetWithDefault("k2", []byte("v22"))
	at.Nil(err)
	at.Equal("v22", string(b2))
}

func Test_Cache_Gorm_Many(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b, err := s.Many([]string{"k1", "k2", "k3"})
	at.Nil(err)
	at.Len(b, 3)
	at.Equal("v1", string(b[0]))
	at.Nil(b[1])
	at.Nil(b[2])
}

func Test_Cache_Gorm_Set(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	err := s.Set("k1", []byte("v11"), time.Minute)
	at.Nil(err)

	var e gormEntry
	at.Nil(s.db.Find(&e, "key = ?", "k1").Error)
	at.Equal("v11", string(e.Value))
}

func Test_Cache_Gorm_Pull(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.Pull("k1")
	at.Nil(err)
	at.Equal("v1", string(b1))
	var e gormEntry
	at.Equal(gorm.ErrRecordNotFound, s.db.First(&e, "key = ?", "k1").Error)

	b2, err := s.Pull("k2")
	at.Nil(err)
	at.Nil(b2)
	at.Equal(nil, s.db.First(&e, "key = ?", "k2").Error)
}

func Test_Cache_Gorm_PullWithDefault(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.PullWithDefault("k1", []byte("v11"))
	at.Nil(err)
	at.Equal("v1", string(b1))
	var e gormEntry
	at.Equal(gorm.ErrRecordNotFound, s.db.First(&e, "key = ?", "k1").Error)

	b2, err := s.PullWithDefault("k2", []byte("v22"))
	at.Nil(err)
	at.Equal("v22", string(b2))
	at.Equal(nil, s.db.First(&e, "key = ?", "k2").Error)
}

func Test_Cache_Gorm_Forever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	err := s.Forever("k1", []byte("v11"))
	at.Nil(err)

	var e gormEntry
	at.Nil(s.db.Find(&e, "key = ?", "k1").Error)
	at.Equal("v11", string(e.Value))
	at.Equal(int64(0), e.Expiry)
}

func Test_Cache_Gorm_Remember(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.Remember("k1", time.Minute, func() ([]byte, error) {
		return []byte("v11"), nil
	})
	at.Nil(err)
	at.Equal("v1", string(b1))

	b2, err := s.Remember("k2", time.Minute, func() ([]byte, error) {
		return []byte("v22"), nil
	})
	at.Nil(err)
	at.Equal("v22", string(b2))
	var e gormEntry
	at.Equal(nil, s.db.First(&e, "key = ?", "k2").Error)
	at.Equal("v22", string(e.Value))
	at.InDelta(time.Now().Add(time.Minute).Unix(), e.Expiry, 1)
}

func Test_Cache_Gorm_RememberForever(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	b1, err := s.RememberForever("k1", func() ([]byte, error) {
		return []byte("v11"), nil
	})
	at.Nil(err)
	at.Equal("v1", string(b1))

	b2, err := s.RememberForever("k2", func() ([]byte, error) {
		return []byte("v22"), nil
	})
	at.Nil(err)
	at.Equal("v22", string(b2))
	var e gormEntry
	at.Equal(nil, s.db.First(&e, "key = ?", "k2").Error)
	at.Equal("v22", string(e.Value))
	at.InDelta(0, e.Expiry, 1)
}

func Test_Cache_Gorm_Delete(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	at.Nil(s.Delete("k1"))
	at.Nil(s.Delete("k2"))
	at.Nil(s.Delete("k3"))

	var e gormEntry
	at.Equal(gorm.ErrRecordNotFound, s.db.First(&e, "key = ?", "k1").Error)
	at.Equal(gorm.ErrRecordNotFound, s.db.First(&e, "key = ?", "k2").Error)
}

func Test_Cache_Gorm_Reset(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	at.Nil(s.Reset())

	var e gormEntry
	at.Equal(gorm.ErrRecordNotFound, s.db.First(&e, "key = ?", "k1").Error)
	at.Equal(gorm.ErrRecordNotFound, s.db.First(&e, "key = ?", "k2").Error)
}

func Test_Cache_Gorm_Close(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	s := getGormStorage(t)

	at.Nil(s.Close())

	at.Panics(func() {
		close(s.done)
	})
}

func Test_Cache_Gorm_GC(t *testing.T) {
	t.Parallel()

	s := getGormStorage(t)

	go s.gc()

	time.Sleep(time.Millisecond * 15)

	close(s.done)

	assert.Eventually(t, func() bool {
		var e gormEntry
		b1 := s.db.First(&e, "key = ?", "k1").Error != gorm.ErrRecordNotFound
		b2 := s.db.First(&e, "key = ?", "k2").Error == gorm.ErrRecordNotFound
		return b1 && b2
	}, time.Second, time.Millisecond*10)
}

type fake struct {
	Entry gormEntry `gorm:"embedded"`
}

func (fake) TableName() string { return "test_dawn_cache" }

func getGormStorage(t *testing.T) *gormStorage {
	s := &gormStorage{
		db:         deck.SetupGormDB(t, &fake{}),
		table:      "test_dawn_cache",
		gcInterval: time.Millisecond * 10,
		done:       make(chan struct{}),
	}

	s.setup()

	// cached k1
	s.db.Create(&gormEntry{
		Key:    "k1",
		Value:  []byte("v1"),
		Expiry: time.Now().Add(time.Minute).Unix(),
	})

	// expired k2
	s.db.Create(&gormEntry{
		Key:    "k2",
		Value:  []byte("v2"),
		Expiry: time.Now().Add(-time.Minute).Unix(),
	})

	return s
}
