package auth

import (
	"errors"
	"testing"

	"github.com/go-dawn/module/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Auth_Service_RegisterByPassword(t *testing.T) {
	at := assert.New(t)

	s, mockRepo, _ := getService()
	var (
		username = "username"
		pass     = "pass"
	)
	mockRepo.On("RegisterByPassword", username, pass).
		Once().Return(1, nil)

	id, err := s.RegisterByPassword(username, pass)

	at.Nil(err)
	at.Equal(1, id)
}

func Test_Auth_Service_RegisterByMobileCode(t *testing.T) {
	at := assert.New(t)

	s, mockRepo, mockValidator := getService()
	var (
		mobile  = "13600008888"
		code    = "123456"
		mockErr = errors.New("fake error")
	)

	t.Run("wrong code", func(t *testing.T) {
		mockValidator.On("Validate", mobile, code).
			Once().Return(mockErr)

		_, err := s.RegisterByMobileCode(mobile, code)

		at.Equal(mockErr, err)
	})

	t.Run("failed", func(t *testing.T) {
		mockValidator.On("Validate", mobile, code).
			Once().Return(nil)
		mockRepo.On("RegisterByMobile", mobile).
			Once().Return(0, mockErr)

		_, err := s.RegisterByMobileCode(mobile, code)

		at.Equal(mockErr, err)
	})

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", mobile, code).
			Once().Return(nil)
		mockRepo.On("RegisterByMobile", mobile).
			Once().Return(1, nil)

		id, err := s.RegisterByMobileCode(mobile, code)

		at.Nil(err)
		at.Equal(1, id)
	})
}

func Test_Auth_Service_RegisterByEmailCode(t *testing.T) {
	at := assert.New(t)

	s, mockRepo, mockValidator := getService()
	var (
		email   = "kiyonlin@gmail.com"
		code    = "123456"
		mockErr = errors.New("fake error")
	)

	t.Run("wrong code", func(t *testing.T) {
		mockValidator.On("Validate", email, code).
			Once().Return(mockErr)

		_, err := s.RegisterByEmailCode(email, code)

		at.Equal(mockErr, err)
	})

	t.Run("failed", func(t *testing.T) {
		mockValidator.On("Validate", email, code).
			Once().Return(nil)
		mockRepo.On("RegisterByEmail", email).
			Once().Return(0, mockErr)

		_, err := s.RegisterByEmailCode(email, code)

		at.Equal(mockErr, err)
	})

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", email, code).
			Once().Return(nil)
		mockRepo.On("RegisterByEmail", email).
			Once().Return(1, nil)

		id, err := s.RegisterByEmailCode(email, code)

		at.Nil(err)
		at.Equal(1, id)
	})
}

func Test_Auth_Service_LoginByPassword(t *testing.T) {
	at := assert.New(t)

	s, mockRepo, _ := getService()
	var (
		username = "username"
		pass     = "pass"
	)
	mockRepo.On("LoginByPassword", username, pass).
		Once().Return(1, nil)

	id, err := s.LoginByPassword(username, pass)

	at.Nil(err)
	at.Equal(1, id)
}

func Test_Auth_Service_LoginByMobileCode(t *testing.T) {
	at := assert.New(t)

	s, mockRepo, mockValidator := getService()
	var (
		mobile  = "13600008888"
		code    = "123456"
		mockErr = errors.New("fake error")
	)

	t.Run("wrong code", func(t *testing.T) {
		mockValidator.On("Validate", mobile, code).
			Once().Return(mockErr)

		_, err := s.LoginByMobileCode(mobile, code)

		at.Equal(mockErr, err)
	})

	t.Run("failed", func(t *testing.T) {
		mockValidator.On("Validate", mobile, code).
			Once().Return(nil)
		mockRepo.On("LoginByMobile", mobile).
			Once().Return(0, mockErr)

		_, err := s.LoginByMobileCode(mobile, code)

		at.Equal(mockErr, err)
	})

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", mobile, code).
			Once().Return(nil)
		mockRepo.On("LoginByMobile", mobile).
			Once().Return(1, nil)

		id, err := s.LoginByMobileCode(mobile, code)

		at.Nil(err)
		at.Equal(1, id)
	})
}

func Test_Auth_Service_LoginByEmailCode(t *testing.T) {
	at := assert.New(t)

	s, mockRepo, mockValidator := getService()
	var (
		email   = "kiyonlin@gmail.com"
		code    = "123456"
		mockErr = errors.New("fake error")
	)

	t.Run("wrong code", func(t *testing.T) {
		mockValidator.On("Validate", email, code).
			Once().Return(mockErr)

		_, err := s.LoginByEmailCode(email, code)

		at.Equal(mockErr, err)
	})

	t.Run("failed", func(t *testing.T) {
		mockValidator.On("Validate", email, code).
			Once().Return(nil)
		mockRepo.On("LoginByEmail", email).
			Once().Return(0, mockErr)

		_, err := s.LoginByEmailCode(email, code)

		at.Equal(mockErr, err)
	})

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", email, code).
			Once().Return(nil)
		mockRepo.On("LoginByEmail", email).
			Once().Return(1, nil)

		id, err := s.LoginByEmailCode(email, code)

		at.Nil(err)
		at.Equal(1, id)
	})
}

func getService() (service, *mocks.Repo, *mocks.CodeValidator) {
	repo, v := new(mocks.Repo), new(mocks.CodeValidator)
	return service{repo: repo, v: v}, repo, v
}
