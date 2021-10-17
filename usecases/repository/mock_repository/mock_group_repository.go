// Code generated by MockGen. DO NOT EDIT.
// Source: group_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/domain"
)

// MockGroupRepository is a mock of GroupRepository interface.
type MockGroupRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGroupRepositoryMockRecorder
}

// MockGroupRepositoryMockRecorder is the mock recorder for MockGroupRepository.
type MockGroupRepositoryMockRecorder struct {
	mock *MockGroupRepository
}

// NewMockGroupRepository creates a new mock instance.
func NewMockGroupRepository(ctrl *gomock.Controller) *MockGroupRepository {
	mock := &MockGroupRepository{ctrl: ctrl}
	mock.recorder = &MockGroupRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupRepository) EXPECT() *MockGroupRepositoryMockRecorder {
	return m.recorder
}

// GetAllGroups mocks base method.
func (m *MockGroupRepository) GetAllGroups() ([]*domain.Groups, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllGroups")
	ret0, _ := ret[0].([]*domain.Groups)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGroups indicates an expected call of GetAllGroups.
func (mr *MockGroupRepositoryMockRecorder) GetAllGroups() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGroups", reflect.TypeOf((*MockGroupRepository)(nil).GetAllGroups))
}

// GetGroup mocks base method.
func (m *MockGroupRepository) GetGroup(groupID uuid.UUID) (*domain.GroupDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", groupID)
	ret0, _ := ret[0].(*domain.GroupDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockGroupRepositoryMockRecorder) GetGroup(groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockGroupRepository)(nil).GetGroup), groupID)
}

// GetGroupsByID mocks base method.
func (m *MockGroupRepository) GetGroupsByID(userID uuid.UUID) ([]*domain.GroupUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsByID", userID)
	ret0, _ := ret[0].([]*domain.GroupUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsByID indicates an expected call of GetGroupsByID.
func (mr *MockGroupRepositoryMockRecorder) GetGroupsByID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsByID", reflect.TypeOf((*MockGroupRepository)(nil).GetGroupsByID), userID)
}
