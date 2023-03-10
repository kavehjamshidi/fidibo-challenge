// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kavehjamshidi/fidibo-challenge/domain"
	mock "github.com/stretchr/testify/mock"
)

// Cacher is an autogenerated mock type for the Cacher type
type Cacher struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *Cacher) Get(_a0 context.Context, _a1 string) (domain.SearchResult, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.SearchResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.SearchResult, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.SearchResult); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.SearchResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0, _a1, _a2
func (_m *Cacher) Store(_a0 context.Context, _a1 string, _a2 domain.SearchResult) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.SearchResult) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCacher interface {
	mock.TestingT
	Cleanup(func())
}

// NewCacher creates a new instance of Cacher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCacher(t mockConstructorTestingTNewCacher) *Cacher {
	mock := &Cacher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
