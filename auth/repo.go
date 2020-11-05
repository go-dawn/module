package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Repo is the repository interface for auth behaviors
type Repo interface {
	// LoginByPassword login system by username and password
	// and return user id if authentication success
	LoginByPassword(username, pass string) (int, error)

	// LoginByMobile login system by mobile number and validate code
	// and return user id if authentication success
	LoginByMobile(mobile, code string) (int, error)

	// LoginByEmail login system by email address and validate code
	// and return user id if authentication success
	LoginByEmail(email, code string) (int, error)
}

// repository is an internal implement of Repo interface
type repository struct {
	db *gorm.DB
}

func defaultRepo(db *gorm.DB) repository {
	return repository{db: db}
}

func (r repository) LoginByPassword(username, pass string) (id int, err error) {
	var u user
	if err = r.db.First(&u, "username = ?", username).Error; err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword(u.Password, []byte(pass)); err != nil {
		return
	}

	id = int(u.ID)

	return
}

func (r repository) LoginByMobile(mobile, code string) (int, error) {
	panic("implement me")
}

func (r repository) LoginByEmail(email, code string) (int, error) {
	panic("implement me")
}

type user struct {
	gorm.Model

	Username string
	Password []byte
}
