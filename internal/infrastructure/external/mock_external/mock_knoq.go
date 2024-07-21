// Code generated by MockGen. DO NOT EDIT.
// Source: knoq.go

// Package mock_external is a generated GoMock package.
package mock_external

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	external "github.com/traPtitech/traPortfolio/internal/infrastructure/external"
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

// GetEvent mocks base method.
func (m *MockKnoqAPI) GetEvent(eventID uuid.UUID) (*external.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", eventID)
	ret0, _ := ret[0].(*external.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockKnoqAPIMockRecorder) GetEvent(eventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockKnoqAPI)(nil).GetEvent), eventID)
}

// GetEvents mocks base method.
func (m *MockKnoqAPI) GetEvents() ([]*external.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents")
	ret0, _ := ret[0].([]*external.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockKnoqAPIMockRecorder) GetEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockKnoqAPI)(nil).GetEvents))
}

// GetEventsByUserID mocks base method.
func (m *MockKnoqAPI) GetEventsByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventsByUserID", userID)
	ret0, _ := ret[0].([]*external.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByUserID indicates an expected call of GetEventsByUserID.
func (mr *MockKnoqAPIMockRecorder) GetEventsByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByUserID", reflect.TypeOf((*MockKnoqAPI)(nil).GetEventsByUserID), userID)
}