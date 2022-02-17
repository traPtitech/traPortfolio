// Code generated by MockGen. DO NOT EDIT.
// Source: knoq.go

// Package mock_external is a generated GoMock package.
package mock_external

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	external "github.com/traPtitech/traPortfolio/interfaces/external"
)

// MockKnoqAPI is a mock of KnoqAPI interface.
type MockKnoqAPI struct {
	ctrl     *gomock.Controller
	recorder *MockKnoqAPIMockRecorder
}

// MockKnoqAPIMockRecorder is the mock recorder for MockKnoqAPI.
type MockKnoqAPIMockRecorder struct {
	mock *MockKnoqAPI
}

// NewMockKnoqAPI creates a new mock instance.
func NewMockKnoqAPI(ctrl *gomock.Controller) *MockKnoqAPI {
	mock := &MockKnoqAPI{ctrl: ctrl}
	mock.recorder = &MockKnoqAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKnoqAPI) EXPECT() *MockKnoqAPIMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockKnoqAPI) GetAll() ([]*external.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*external.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockKnoqAPIMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockKnoqAPI)(nil).GetAll))
}

// GetByID mocks base method.
func (m *MockKnoqAPI) GetByID(id uuid.UUID) (*external.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*external.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockKnoqAPIMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockKnoqAPI)(nil).GetByID), id)
}

// GetByUserID mocks base method.
func (m *MockKnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", userID)
	ret0, _ := ret[0].([]*external.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockKnoqAPIMockRecorder) GetByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockKnoqAPI)(nil).GetByUserID), userID)
}
