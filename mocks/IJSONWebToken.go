// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	token "github.com/aasumitro/tix/pkg/token"
	mock "github.com/stretchr/testify/mock"
)

// IJSONWebToken is an autogenerated mock type for the IJSONWebToken type
type IJSONWebToken struct {
	mock.Mock
}

// Claim provides a mock function with given fields: payload
func (_m *IJSONWebToken) Claim(payload interface{}) (string, error) {
	ret := _m.Called(payload)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}) (string, error)); ok {
		return rf(payload)
	}
	if rf, ok := ret.Get(0).(func(interface{}) string); ok {
		r0 = rf(payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtractAndValidateJWT provides a mock function with given fields: secret, cookie
func (_m *IJSONWebToken) ExtractAndValidateJWT(secret string, cookie *http.Cookie) (*token.JSONWebTokenClaim, error) {
	ret := _m.Called(secret, cookie)

	var r0 *token.JSONWebTokenClaim
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *http.Cookie) (*token.JSONWebTokenClaim, error)); ok {
		return rf(secret, cookie)
	}
	if rf, ok := ret.Get(0).(func(string, *http.Cookie) *token.JSONWebTokenClaim); ok {
		r0 = rf(secret, cookie)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token.JSONWebTokenClaim)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *http.Cookie) error); ok {
		r1 = rf(secret, cookie)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIJSONWebToken interface {
	mock.TestingT
	Cleanup(func())
}

// NewIJSONWebToken creates a new instance of IJSONWebToken. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIJSONWebToken(t mockConstructorTestingTNewIJSONWebToken) *IJSONWebToken {
	mock := &IJSONWebToken{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}