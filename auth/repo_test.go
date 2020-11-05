package auth

import (
	"testing"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"

	"github.com/go-dawn/pkg/deck"
)

func Test_Auth_Repo_LoginByPassword(t *testing.T) {
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

func getRepo(t *testing.T) repository {
	gdb := deck.SetupGormDB(t, &user{})
	return repository{gdb}
}

func (r repository) createUser(t *testing.T, username, pass string) *user {
	at := assert.New(t)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	at.Nil(err)

	u := &user{Username: username, Password: hashedPass}

	at.Nil(r.db.Create(u).Error)

	return u
}
