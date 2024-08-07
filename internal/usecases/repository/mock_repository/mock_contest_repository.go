// Code generated by MockGen. DO NOT EDIT.
// Source: contest_repository.go

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

// MockContestRepository is a mock of ContestRepository interface.
type MockContestRepository struct {
	ctrl     *gomock.Controller
	recorder *MockContestRepositoryMockRecorder
}

// MockContestRepositoryMockRecorder is the mock recorder for MockContestRepository.
type MockContestRepositoryMockRecorder struct {
	mock *MockContestRepository
}

// NewMockContestRepository creates a new mock instance.
func NewMockContestRepository(ctrl *gomock.Controller) *MockContestRepository {
	mock := &MockContestRepository{ctrl: ctrl}
	mock.recorder = &MockContestRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContestRepository) EXPECT() *MockContestRepositoryMockRecorder {
	return m.recorder
}

// CreateContest mocks base method.
func (m *MockContestRepository) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.ContestDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContest", ctx, args)
	ret0, _ := ret[0].(*domain.ContestDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContest indicates an expected call of CreateContest.
func (mr *MockContestRepositoryMockRecorder) CreateContest(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContest", reflect.TypeOf((*MockContestRepository)(nil).CreateContest), ctx, args)
}

// CreateContestTeam mocks base method.
func (m *MockContestRepository) CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContestTeam", ctx, contestID, args)
	ret0, _ := ret[0].(*domain.ContestTeamDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContestTeam indicates an expected call of CreateContestTeam.
func (mr *MockContestRepositoryMockRecorder) CreateContestTeam(ctx, contestID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContestTeam", reflect.TypeOf((*MockContestRepository)(nil).CreateContestTeam), ctx, contestID, args)
}

// DeleteContest mocks base method.
func (m *MockContestRepository) DeleteContest(ctx context.Context, contestID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContest", ctx, contestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContest indicates an expected call of DeleteContest.
func (mr *MockContestRepositoryMockRecorder) DeleteContest(ctx, contestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContest", reflect.TypeOf((*MockContestRepository)(nil).DeleteContest), ctx, contestID)
}

// DeleteContestTeam mocks base method.
func (m *MockContestRepository) DeleteContestTeam(ctx context.Context, contestID, teamID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContestTeam", ctx, contestID, teamID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContestTeam indicates an expected call of DeleteContestTeam.
func (mr *MockContestRepositoryMockRecorder) DeleteContestTeam(ctx, contestID, teamID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContestTeam", reflect.TypeOf((*MockContestRepository)(nil).DeleteContestTeam), ctx, contestID, teamID)
}

// EditContestTeamMembers mocks base method.
func (m *MockContestRepository) EditContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditContestTeamMembers", ctx, teamID, memberIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditContestTeamMembers indicates an expected call of EditContestTeamMembers.
func (mr *MockContestRepositoryMockRecorder) EditContestTeamMembers(ctx, teamID, memberIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditContestTeamMembers", reflect.TypeOf((*MockContestRepository)(nil).EditContestTeamMembers), ctx, teamID, memberIDs)
}

// GetContest mocks base method.
func (m *MockContestRepository) GetContest(ctx context.Context, contestID uuid.UUID) (*domain.ContestDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContest", ctx, contestID)
	ret0, _ := ret[0].(*domain.ContestDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContest indicates an expected call of GetContest.
func (mr *MockContestRepositoryMockRecorder) GetContest(ctx, contestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContest", reflect.TypeOf((*MockContestRepository)(nil).GetContest), ctx, contestID)
}

// GetContestTeam mocks base method.
func (m *MockContestRepository) GetContestTeam(ctx context.Context, contestID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContestTeam", ctx, contestID, teamID)
	ret0, _ := ret[0].(*domain.ContestTeamDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContestTeam indicates an expected call of GetContestTeam.
func (mr *MockContestRepositoryMockRecorder) GetContestTeam(ctx, contestID, teamID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContestTeam", reflect.TypeOf((*MockContestRepository)(nil).GetContestTeam), ctx, contestID, teamID)
}

// GetContestTeamMembers mocks base method.
func (m *MockContestRepository) GetContestTeamMembers(ctx context.Context, contestID, teamID uuid.UUID) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContestTeamMembers", ctx, contestID, teamID)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContestTeamMembers indicates an expected call of GetContestTeamMembers.
func (mr *MockContestRepositoryMockRecorder) GetContestTeamMembers(ctx, contestID, teamID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContestTeamMembers", reflect.TypeOf((*MockContestRepository)(nil).GetContestTeamMembers), ctx, contestID, teamID)
}

// GetContestTeams mocks base method.
func (m *MockContestRepository) GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContestTeams", ctx, contestID)
	ret0, _ := ret[0].([]*domain.ContestTeam)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContestTeams indicates an expected call of GetContestTeams.
func (mr *MockContestRepositoryMockRecorder) GetContestTeams(ctx, contestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContestTeams", reflect.TypeOf((*MockContestRepository)(nil).GetContestTeams), ctx, contestID)
}

// GetContests mocks base method.
func (m *MockContestRepository) GetContests(ctx context.Context) ([]*domain.Contest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContests", ctx)
	ret0, _ := ret[0].([]*domain.Contest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContests indicates an expected call of GetContests.
func (mr *MockContestRepositoryMockRecorder) GetContests(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContests", reflect.TypeOf((*MockContestRepository)(nil).GetContests), ctx)
}

// UpdateContest mocks base method.
func (m *MockContestRepository) UpdateContest(ctx context.Context, contestID uuid.UUID, args *repository.UpdateContestArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContest", ctx, contestID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContest indicates an expected call of UpdateContest.
func (mr *MockContestRepositoryMockRecorder) UpdateContest(ctx, contestID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContest", reflect.TypeOf((*MockContestRepository)(nil).UpdateContest), ctx, contestID, args)
}

// UpdateContestTeam mocks base method.
func (m *MockContestRepository) UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContestTeam", ctx, teamID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContestTeam indicates an expected call of UpdateContestTeam.
func (mr *MockContestRepositoryMockRecorder) UpdateContestTeam(ctx, teamID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContestTeam", reflect.TypeOf((*MockContestRepository)(nil).UpdateContestTeam), ctx, teamID, args)
}
