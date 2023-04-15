// Code generated by mockery v2.22.1. DO NOT EDIT.

package battle

import (
	vo "github.com/Amobe/PlayGame/server/internal/domain/vo"
	mock "github.com/stretchr/testify/mock"
)

// MockUnit is an autogenerated mock type for the Unit type
type MockUnit struct {
	mock.Mock
}

// GetAgi provides a mock function with given fields:
func (_m *MockUnit) GetAgi() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetAttributeMap provides a mock function with given fields:
func (_m *MockUnit) GetAttributeMap() vo.AttributeMap {
	ret := _m.Called()

	var r0 vo.AttributeMap
	if rf, ok := ret.Get(0).(func() vo.AttributeMap); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(vo.AttributeMap)
		}
	}

	return r0
}

// GetGroundIdx provides a mock function with given fields:
func (_m *MockUnit) GetGroundIdx() vo.GroundIdx {
	ret := _m.Called()

	var r0 vo.GroundIdx
	if rf, ok := ret.Get(0).(func() vo.GroundIdx); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(vo.GroundIdx)
	}

	return r0
}

// GetSkill provides a mock function with given fields:
func (_m *MockUnit) GetSkill() vo.Skill {
	ret := _m.Called()

	var r0 vo.Skill
	if rf, ok := ret.Get(0).(func() vo.Skill); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(vo.Skill)
	}

	return r0
}

// IsDead provides a mock function with given fields:
func (_m *MockUnit) IsDead() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TakeAffect provides a mock function with given fields: affects
func (_m *MockUnit) TakeAffect(affects vo.AttributeMap) vo.GroundUnit {
	ret := _m.Called(affects)

	var r0 vo.GroundUnit
	if rf, ok := ret.Get(0).(func(vo.AttributeMap) vo.GroundUnit); ok {
		r0 = rf(affects)
	} else {
		r0 = ret.Get(0).(vo.GroundUnit)
	}

	return r0
}

type mockConstructorTestingTNewMockUnit interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockUnit creates a new instance of MockUnit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUnit(t mockConstructorTestingTNewMockUnit) *MockUnit {
	mock := &MockUnit{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}