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
func (m *MockUserRepository) CreateAccount(userID uuid.UUID, args *repository.CreateAccountArgs) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", userID, args)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockUserRepositoryMockRecorder) CreateAccount(userID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockUserRepository)(nil).CreateAccount), userID, args)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(args repository.CreateUserArgs) (*domain.UserDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", args)
	ret0, _ := ret[0].(*domain.UserDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), args)
}

// DeleteAccount mocks base method.
func (m *MockUserRepository) DeleteAccount(userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockUserRepositoryMockRecorder) DeleteAccount(userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockUserRepository)(nil).DeleteAccount), userID, accountID)
}

// GetAccount mocks base method.
func (m *MockUserRepository) GetAccount(userID, accountID uuid.UUID) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", userID, accountID)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockUserRepositoryMockRecorder) GetAccount(userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockUserRepository)(nil).GetAccount), userID, accountID)
}

// GetAccounts mocks base method.
func (m *MockUserRepository) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", userID)
	ret0, _ := ret[0].([]*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockUserRepositoryMockRecorder) GetAccounts(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockUserRepository)(nil).GetAccounts), userID)
}

// GetContests mocks base method.
func (m *MockUserRepository) GetContests(userID uuid.UUID) ([]*domain.UserContest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContests", userID)
	ret0, _ := ret[0].([]*domain.UserContest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContests indicates an expected call of GetContests.
func (mr *MockUserRepositoryMockRecorder) GetContests(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContests", reflect.TypeOf((*MockUserRepository)(nil).GetContests), userID)
}

// GetGroupsByUserID mocks base method.
func (m *MockUserRepository) GetGroupsByUserID(userID uuid.UUID) ([]*domain.GroupUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsByUserID", userID)
	ret0, _ := ret[0].([]*domain.GroupUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsByUserID indicates an expected call of GetGroupsByUserID.
func (mr *MockUserRepositoryMockRecorder) GetGroupsByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsByUserID", reflect.TypeOf((*MockUserRepository)(nil).GetGroupsByUserID), userID)
}

// GetProjects mocks base method.
func (m *MockUserRepository) GetProjects(userID uuid.UUID) ([]*domain.UserProject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjects", userID)
	ret0, _ := ret[0].([]*domain.UserProject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjects indicates an expected call of GetProjects.
func (mr *MockUserRepositoryMockRecorder) GetProjects(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjects", reflect.TypeOf((*MockUserRepository)(nil).GetProjects), userID)
}

// GetUser mocks base method.
func (m *MockUserRepository) GetUser(userID uuid.UUID) (*domain.UserDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", userID)
	ret0, _ := ret[0].(*domain.UserDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserRepositoryMockRecorder) GetUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserRepository)(nil).GetUser), userID)
}

// GetUsers mocks base method.
func (m *MockUserRepository) GetUsers(args *repository.GetUsersArgs) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", args)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserRepositoryMockRecorder) GetUsers(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserRepository)(nil).GetUsers), args)
}

// UpdateAccount mocks base method.
func (m *MockUserRepository) UpdateAccount(userID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", userID, accountID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockUserRepositoryMockRecorder) UpdateAccount(userID, accountID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockUserRepository)(nil).UpdateAccount), userID, accountID, args)
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(userID uuid.UUID, args *repository.UpdateUserArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", userID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(userID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), userID, args)
}
