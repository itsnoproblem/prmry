// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/itsnoproblem/prmry/pkg/interaction"

	mock "github.com/stretchr/testify/mock"
)

// ModerationRepo is an autogenerated mock type for the ModerationRepo type
type ModerationRepo struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, mod
func (_m *ModerationRepo) Add(ctx context.Context, mod interaction.Moderation) error {
	ret := _m.Called(ctx, mod)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interaction.Moderation) error); ok {
		r0 = rf(ctx, mod)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// All provides a mock function with given fields: ctx
func (_m *ModerationRepo) All(ctx context.Context) ([]interaction.Moderation, error) {
	ret := _m.Called(ctx)

	var r0 []interaction.Moderation
	if rf, ok := ret.Get(0).(func(context.Context) []interaction.Moderation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interaction.Moderation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: ctx, id
func (_m *ModerationRepo) Remove(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewModerationRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewModerationRepo creates a new instance of ModerationRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewModerationRepo(t mockConstructorTestingTNewModerationRepo) *ModerationRepo {
	mock := &ModerationRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
