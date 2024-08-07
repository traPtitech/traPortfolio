// Code generated by MockGen. DO NOT EDIT.
// Source: traq.go

// Package mock_external is a generated GoMock package.
package mock_external

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	external "github.com/traPtitech/traPortfolio/internal/infrastructure/external"
)

// MockTraQAPI is a mock of TraQAPI interface.
type MockTraQAPI struct {
	ctrl     *gomock.Controller
	recorder *MockTraQAPIMockRecorder
}

// MockTraQAPIMockRecorder is the mock recorder for MockTraQAPI.
type MockTraQAPIMockRecorder struct {
	mock *MockTraQAPI
}

// NewMockTraQAPI creates a new mock instance.
func NewMockTraQAPI(ctrl *gomock.Controller) *MockTraQAPI {
	mock := &MockTraQAPI{ctrl: ctrl}
	mock.recorder = &MockTraQAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTraQAPI) EXPECT() *MockTraQAPIMockRecorder {
	return m.recorder
}

// GetUsers mocks base method.
func (m *MockTraQAPI) GetUsers(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", args)
	ret0, _ := ret[0].([]*external.TraQUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockTraQAPIMockRecorder) GetUsers(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockTraQAPI)(nil).GetUsers), args)
}
