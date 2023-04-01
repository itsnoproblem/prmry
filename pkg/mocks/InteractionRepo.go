// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interaction"

	mock "github.com/stretchr/testify/mock"
)

// InteractionRepo is an autogenerated mock type for the InteractionRepo type
type InteractionRepo struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, in
func (_m *InteractionRepo) Add(ctx context.Context, in interaction.Interaction) (string, error) {
	ret := _m.Called(ctx, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, interaction.Interaction) string); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interaction.Interaction) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Interaction provides a mock function with given fields: ctx, id
func (_m *InteractionRepo) Interaction(ctx context.Context, id string) (interaction.Interaction, error) {
	ret := _m.Called(ctx, id)

	var r0 interaction.Interaction
	if rf, ok := ret.Get(0).(func(context.Context, string) interaction.Interaction); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(interaction.Interaction)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: ctx, id
func (_m *InteractionRepo) Remove(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Summaries provides a mock function with given fields: ctx
func (_m *InteractionRepo) Summaries(ctx context.Context) ([]interaction.Summary, error) {
	ret := _m.Called(ctx)

	var r0 []interaction.Summary
	if rf, ok := ret.Get(0).(func(context.Context) []interaction.Summary); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interaction.Summary)
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

type mockConstructorTestingTNewInteractionRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewInteractionRepo creates a new instance of InteractionRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInteractionRepo(t mockConstructorTestingTNewInteractionRepo) *InteractionRepo {
	mock := &InteractionRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
