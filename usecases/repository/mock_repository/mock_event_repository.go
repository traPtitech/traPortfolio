// Code generated by MockGen. DO NOT EDIT.
// Source: event_repository.go

// Package mock_repository isb a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/domain"
	model "github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

// MockEventRepository is a mock of EventRepository interface.
type MockEventRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEventRepositoryMockRecorder
}

// MockEventRepositoryMockRecorder is the mock recorder for MockEventRepository.
type MockEventRepositoryMockRecorder struct {
	mock *MockEventRepository
}

// NewMockEventRepository creates a new mock instance.
func NewMockEventRepository(ctrl *gomock.Controller) *MockEventRepository {
	mock := &MockEventRepository{ctrl: ctrl}
	mock.recorder = &MockEventRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventRepository) EXPECT() *MockEventRepositoryMockRecorder {
	return m.recorder
}

// GetEvent mocks base method.
func (m *MockEventRepository) GetEvent(id uuid.UUID) (*domain.EventDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", id)
	ret0, _ := ret[0].(*domain.EventDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockEventRepositoryMockRecorder) GetEvent(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockEventRepository)(nil).GetEvent), id)
}

// GetEvents mocks base method.
func (m *MockEventRepository) GetEvents() ([]*domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents")
	ret0, _ := ret[0].([]*domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockEventRepositoryMockRecorder) GetEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockEventRepository)(nil).GetEvents))
}

// GetUserEvents mocks base method.
func (m *MockEventRepository) GetUserEvents(userID uuid.UUID) ([]*domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserEvents", userID)
	ret0, _ := ret[0].([]*domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserEvents indicates an expected call of GetUserEvents.
func (mr *MockEventRepositoryMockRecorder) GetUserEvents(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserEvents", reflect.TypeOf((*MockEventRepository)(nil).GetUserEvents), userID)
}

// UpdateEvent mocks base method.
func (m *MockEventRepository) UpdateEvent(elv *model.EventLevelRelation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", elv)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventRepositoryMockRecorder) UpdateEvent(elv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventRepository)(nil).UpdateEvent), elv)
}
