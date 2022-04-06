// Code generated by mockery v2.10.0. DO NOT EDIT.

package advert

import (
	context "context"

	advert "github.com/ukrainian-brothers/board-backend/domain/advert"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// RepositoryMock is an autogenerated mock type for the RepositoryMock type
type RepositoryMock struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, _a1
func (_m *RepositoryMock) Add(ctx context.Context, _a1 *advert.Advert) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *advert.Advert) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *RepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *RepositoryMock) Get(ctx context.Context, id uuid.UUID) (advert.Advert, error) {
	ret := _m.Called(ctx, id)

	var r0 advert.Advert
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) advert.Advert); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(advert.Advert)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}