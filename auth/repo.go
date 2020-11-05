package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var bcryptCost = bcrypt.DefaultCost

// Repo is the repository interface for auth behaviors
type Repo interface {
	// RegisterByPassword gets a new user by username and password
	RegisterByPassword(username, pass string) (int, error)

	// RegisterByMobileCode gets a new user by mobile
	RegisterByMobile(mobile string) (int, error)

	// RegisterByEmailCode gets a new user by email
	RegisterByEmail(email string) (int, error)

	// LoginByPassword login system by username and password
	// and return user id if authentication success
	LoginByPassword(username, pass string) (int, error)

	// LoginByMobileCode login system by mobile number
	// and return user id if authentication success
	LoginByMobile(mobile string) (int, error)

	// LoginByEmailCode login system by email address
	// and return user id if authentication success
	LoginByEmail(email string) (int, error)
}

// repository is an internal implement of Repo interface
type repository struct {
	db *gorm.DB
}

func (r repository) RegisterByPassword(username, pass string) (id int, err error) {
	u := &user{Username: username}

	if u.Password, err = bcrypt.GenerateFromPassword([]byte(pass), bcryptCost); err != nil {
		return
	}

	err, id = r.db.Create(u).Error, int(u.ID)

	return
}

func (r repository) RegisterByMobile(mobile string) (int, error) {
	u := &user{Mobile: mobile}
	err, id := r.db.Create(u).Error, int(u.ID)
	return id, err
}

func (r repository) RegisterByEmail(email string) (int, error) {
	u := &user{Email: email}
	err, id := r.db.Create(u).Error, int(u.ID)
	return id, err
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

func (r repository) LoginByMobile(mobile string) (int, error) {
	var u user
	err := r.db.First(&u, "mobile = ?", mobile).Error
	return int(u.ID), err
}

func (r repository) LoginByEmail(email string) (int, error) {
	var u user
	err := r.db.First(&u, "email = ?", email).Error
	return int(u.ID), err
}

type user struct {
	gorm.Model

	Username string `gorm:"uniqueIndex"`
	Password []byte
	Mobile   string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
}
