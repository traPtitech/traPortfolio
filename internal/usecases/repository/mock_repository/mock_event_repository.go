// Code generated by MockGen. DO NOT EDIT.
// Source: event_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/internal/domain"
	repository "github.com/traPtitech/traPortfolio/internal/usecases/repository"
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

// CreateEventLevel mocks base method.
func (m *MockEventRepository) CreateEventLevel(ctx context.Context, args *repository.CreateEventLevelArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEventLevel", ctx, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEventLevel indicates an expected call of CreateEventLevel.
func (mr *MockEventRepositoryMockRecorder) CreateEventLevel(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEventLevel", reflect.TypeOf((*MockEventRepository)(nil).CreateEventLevel), ctx, args)
}

// GetEvent mocks base method.
func (m *MockEventRepository) GetEvent(ctx context.Context, eventID uuid.UUID) (*domain.EventDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", ctx, eventID)
	ret0, _ := ret[0].(*domain.EventDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockEventRepositoryMockRecorder) GetEvent(ctx, eventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockEventRepository)(nil).GetEvent), ctx, eventID)
}

// GetEvents mocks base method.
func (m *MockEventRepository) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", ctx)
	ret0, _ := ret[0].([]*domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockEventRepositoryMockRecorder) GetEvents(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockEventRepository)(nil).GetEvents), ctx)
}

// GetUserEvents mocks base method.
func (m *MockEventRepository) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserEvents", ctx, userID)
	ret0, _ := ret[0].([]*domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserEvents indicates an expected call of GetUserEvents.
func (mr *MockEventRepositoryMockRecorder) GetUserEvents(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserEvents", reflect.TypeOf((*MockEventRepository)(nil).GetUserEvents), ctx, userID)
}

// UpdateEventLevel mocks base method.
func (m *MockEventRepository) UpdateEventLevel(ctx context.Context, eventID uuid.UUID, args *repository.UpdateEventLevelArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEventLevel", ctx, eventID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEventLevel indicates an expected call of UpdateEventLevel.
func (mr *MockEventRepositoryMockRecorder) UpdateEventLevel(ctx, eventID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEventLevel", reflect.TypeOf((*MockEventRepository)(nil).UpdateEventLevel), ctx, eventID, args)
}
