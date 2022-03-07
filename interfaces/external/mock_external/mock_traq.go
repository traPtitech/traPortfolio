// Code generated by MockGen. DO NOT EDIT.
// Source: traq.go

// Package mock_external is a generated GoMock package.
package mock_external

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	external "github.com/traPtitech/traPortfolio/interfaces/external"
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

// GetByID mocks base method.
func (m *MockTraQAPI) GetByID(id uuid.UUID) (*external.TraQUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*external.TraQUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockTraQAPIMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTraQAPI)(nil).GetByID), id)
}