package auth

import (
	"testing"

	"github.com/go-dawn/pkg/deck"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Test_Auth_Repo_RegisterByPassword(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	var (
		username = "username"
		pass     = "pass"
	)

	t.Run("success", func(t *testing.T) {
		repo := getRepo(t)

		id, err := repo.RegisterByPassword(username, pass)
		at.Nil(err)
		at.Equal(1, id)
	})

	t.Run("exist", func(t *testing.T) {
		repo := getRepo(t)
		repo.createUser(t, username, pass)

		_, err := repo.RegisterByPassword(username, pass)
		at.NotNil(err)
		at.Contains(err.Error(), "username")
	})
}

func Test_Auth_Repo_RegisterByMobile(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	mobile := "13600008888"

	t.Run("success", func(t *testing.T) {
		repo := getRepo(t)
		id, err := repo.RegisterByMobile(mobile)
		at.Nil(err)
		at.Equal(1, id)
	})

	t.Run("exist", func(t *testing.T) {
		repo := getRepo(t)
		repo.createMobileUser(t, mobile)

		_, err := repo.RegisterByMobile(mobile)
		at.NotNil(err)
		at.Contains(err.Error(), "mobile")
	})
}

func Test_Auth_Repo_RegisterByEmail(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	email := "kiyonlin@gmail.com"

	t.Run("success", func(t *testing.T) {
		repo := getRepo(t)
		id, err := repo.RegisterByEmail(email)
		at.Nil(err)
		at.Equal(1, id)
	})

	t.Run("exist", func(t *testing.T) {
		repo := getRepo(t)
		repo.createEmailUser(t, email)

		_, err := repo.RegisterByEmail(email)
		at.NotNil(err)
		at.Contains(err.Error(), "email")
	})
}

func Test_Auth_Repo_LoginByPassword(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	var (
		username = "username"
		pass     = "pass"
	)

	t.Run("non-exist", func(t *testing.T) {
		repo := getRepo(t)

		_, err := repo.LoginByPassword(username, pass)

		at.Equal(gorm.ErrRecordNotFound, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := getRepo(t)

		repo.createUser(t, username, pass)

		_, err := repo.LoginByPassword(username, pass+"1")
		at.Equal(bcrypt.ErrMismatchedHashAndPassword, err)
	})

	t.Run("success", func(t *testing.T) {
		repo := getRepo(t)

		repo.createUser(t, username, pass)

		id, err := repo.LoginByPassword(username, pass)
		at.Nil(err)
		at.Equal(1, id)
	})
}

func Test_Auth_Repo_LoginByMobile(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	mobile := "13600008888"

	t.Run("non-exist", func(t *testing.T) {
		repo := getRepo(t)

		_, err := repo.LoginByMobile(mobile)

		at.Equal(gorm.ErrRecordNotFound, err)
	})

	t.Run("success", func(t *testing.T) {
		repo := getRepo(t)

		repo.createMobileUser(t, mobile)

		id, err := repo.LoginByMobile(mobile)
		at.Nil(err)
		at.Equal(1, id)
	})
}

func Test_Auth_Repo_LoginByEmail(t *testing.T) {
	t.Parallel()

	at := assert.New(t)
	email := "kiyonlin@gmail.com"

	t.Run("non-exist", func(t *testing.T) {
		repo := getRepo(t)

		_, err := repo.LoginByEmail(email)

		at.Equal(gorm.ErrRecordNotFound, err)
	})

	t.Run("success", func(t *testing.T) {
		repo := getRepo(t)

		repo.createEmailUser(t, email)

		id, err := repo.LoginByEmail(email)
		at.Nil(err)
		at.Equal(1, id)
	})
}

func getRepo(t *testing.T) repository {
	gdb := deck.SetupGormDB(t, &user{})
	return repository{gdb}
}

func (r repository) createUser(t *testing.T, username, pass string) *user {
	at := assert.New(t)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcryptCost)
	at.Nil(err)

	u := &user{Username: username, Password: hashedPass}

	at.Nil(r.db.Select("username", "password").Create(u).Error)

	return u
}

func (r repository) createMobileUser(t *testing.T, mobile string) *user {
	at := assert.New(t)

	u := &user{Mobile: mobile}

	at.Nil(r.db.Select("mobile").Create(u).Error)

	return u
}

func (r repository) createEmailUser(t *testing.T, email string) *user {
	at := assert.New(t)

	u := &user{Email: email}

	at.Nil(r.db.Select("email").Create(u).Error)

	return u
}
