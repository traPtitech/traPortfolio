// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/domain"
	repository "github.com/traPtitech/traPortfolio/usecases/repository"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockUserService) CreateAccount(ctx context.Context, userID uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, userID, account)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockUserServiceMockRecorder) CreateAccount(ctx, userID, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockUserService)(nil).CreateAccount), ctx, userID, account)
}

// DeleteAccount mocks base method.
func (m *MockUserService) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockUserServiceMockRecorder) DeleteAccount(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockUserService)(nil).DeleteAccount), ctx, userID, accountID)
}

// EditAccount mocks base method.
func (m *MockUserService) EditAccount(ctx context.Context, userID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditAccount", ctx, userID, accountID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditAccount indicates an expected call of EditAccount.
func (mr *MockUserServiceMockRecorder) EditAccount(ctx, userID, accountID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditAccount", reflect.TypeOf((*MockUserService)(nil).EditAccount), ctx, userID, accountID, args)
}

// GetAccount mocks base method.
func (m *MockUserService) GetAccount(userID, accountID uuid.UUID) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", userID, accountID)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockUserServiceMockRecorder) GetAccount(userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockUserService)(nil).GetAccount), userID, accountID)
}

// GetAccounts mocks base method.
func (m *MockUserService) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", userID)
	ret0, _ := ret[0].([]*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockUserServiceMockRecorder) GetAccounts(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockUserService)(nil).GetAccounts), userID)
}

// GetGroupsByUserID mocks base method.
func (m *MockUserService) GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.GroupUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsByUserID", ctx, userID)
	ret0, _ := ret[0].([]*domain.GroupUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsByUserID indicates an expected call of GetGroupsByUserID.
func (mr *MockUserServiceMockRecorder) GetGroupsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsByUserID", reflect.TypeOf((*MockUserService)(nil).GetGroupsByUserID), ctx, userID)
}

// GetUser mocks base method.
func (m *MockUserService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userID)
	ret0, _ := ret[0].(*domain.UserDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserServiceMockRecorder) GetUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserService)(nil).GetUser), ctx, userID)
}

// GetUserContests mocks base method.
func (m *MockUserService) GetUserContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserContests", ctx, userID)
	ret0, _ := ret[0].([]*domain.UserContest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserContests indicates an expected call of GetUserContests.
func (mr *MockUserServiceMockRecorder) GetUserContests(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserContests", reflect.TypeOf((*MockUserService)(nil).GetUserContests), ctx, userID)
}

// GetUserEvents mocks base method.
func (m *MockUserService) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserEvents", ctx, userID)
	ret0, _ := ret[0].([]*domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserEvents indicates an expected call of GetUserEvents.
func (mr *MockUserServiceMockRecorder) GetUserEvents(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserEvents", reflect.TypeOf((*MockUserService)(nil).GetUserEvents), ctx, userID)
}

// GetUserProjects mocks base method.
func (m *MockUserService) GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProjects", ctx, userID)
	ret0, _ := ret[0].([]*domain.UserProject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProjects indicates an expected call of GetUserProjects.
func (mr *MockUserServiceMockRecorder) GetUserProjects(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProjects", reflect.TypeOf((*MockUserService)(nil).GetUserProjects), ctx, userID)
}

// GetUsers mocks base method.
func (m *MockUserService) GetUsers(ctx context.Context, args *repository.GetUsersArgs) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx, args)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserServiceMockRecorder) GetUsers(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserService)(nil).GetUsers), ctx, args)
}

// Update mocks base method.
func (m *MockUserService) Update(ctx context.Context, userID uuid.UUID, args *repository.UpdateUserArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, userID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserServiceMockRecorder) Update(ctx, userID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserService)(nil).Update), ctx, userID, args)
}
