// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// Delete provides a mock function with given fields: key
func (_m *Storage) Delete(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *Storage) Get(key string) ([]byte, error) {
	ret := _m.Called(key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: key, value, ttl
func (_m *Storage) Set(key string, value []byte, ttl time.Duration) error {
	ret := _m.Called(key, value, ttl)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, time.Duration) error); ok {
		r0 = rf(key, value, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
