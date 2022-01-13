// Code generated by MockGen. DO NOT EDIT.
// Source: user_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/domain"
	repository "github.com/traPtitech/traPortfolio/usecases/repository"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockUserRepository) CreateAccount(id uuid.UUID, args *repository.CreateAccountArgs) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", id, args)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockUserRepositoryMockRecorder) CreateAccount(id, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockUserRepository)(nil).CreateAccount), id, args)
}

// DeleteAccount mocks base method.
func (m *MockUserRepository) DeleteAccount(id, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", id, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockUserRepositoryMockRecorder) DeleteAccount(id, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockUserRepository)(nil).DeleteAccount), id, accountID)
}

// GetAccount mocks base method.
func (m *MockUserRepository) GetAccount(id, accountID uuid.UUID) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", id, accountID)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockUserRepositoryMockRecorder) GetAccount(id, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockUserRepository)(nil).GetAccount), id, accountID)
}

// GetAccounts mocks base method.
func (m *MockUserRepository) GetAccounts(id uuid.UUID) ([]*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", id)
	ret0, _ := ret[0].([]*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockUserRepositoryMockRecorder) GetAccounts(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockUserRepository)(nil).GetAccounts), id)
}

// GetContests mocks base method.
func (m *MockUserRepository) GetContests(id uuid.UUID) ([]*domain.UserContest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContests", id)
	ret0, _ := ret[0].([]*domain.UserContest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContests indicates an expected call of GetContests.
func (mr *MockUserRepositoryMockRecorder) GetContests(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContests", reflect.TypeOf((*MockUserRepository)(nil).GetContests), id)
}

// GetGroupsByUserID mocks base method.
func (m *MockUserRepository) GetGroupsByUserID(id uuid.UUID) ([]*domain.GroupUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsByUserID", id)
	ret0, _ := ret[0].([]*domain.GroupUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsByUserID indicates an expected call of GetGroupsByUserID.
func (mr *MockUserRepositoryMockRecorder) GetGroupsByUserID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsByUserID", reflect.TypeOf((*MockUserRepository)(nil).GetGroupsByUserID), id)
}

// GetProjects mocks base method.
func (m *MockUserRepository) GetProjects(id uuid.UUID) ([]*domain.UserProject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjects", id)
	ret0, _ := ret[0].([]*domain.UserProject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjects indicates an expected call of GetProjects.
func (mr *MockUserRepositoryMockRecorder) GetProjects(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjects", reflect.TypeOf((*MockUserRepository)(nil).GetProjects), id)
}

// GetUser mocks base method.
func (m *MockUserRepository) GetUser(id uuid.UUID) (*domain.UserDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*domain.UserDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserRepositoryMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserRepository)(nil).GetUser), id)
}

// GetUsers mocks base method.
func (m *MockUserRepository) GetUsers() ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers")
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserRepositoryMockRecorder) GetUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserRepository)(nil).GetUsers))
}

// Update mocks base method.
func (m *MockUserRepository) Update(id uuid.UUID, changes map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, changes)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(id, changes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), id, changes)
}

// UpdateAccount mocks base method.
func (m *MockUserRepository) UpdateAccount(id, accountID uuid.UUID, changes map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", id, accountID, changes)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockUserRepositoryMockRecorder) UpdateAccount(id, accountID, changes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockUserRepository)(nil).UpdateAccount), id, accountID, changes)
}
