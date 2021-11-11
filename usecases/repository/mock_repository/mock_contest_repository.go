// Code generated by MockGen. DO NOT EDIT.
// Source: contest_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/domain"
	repository "github.com/traPtitech/traPortfolio/usecases/repository"
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

// AddContestTeamMembers mocks base method.
func (m *MockContestRepository) AddContestTeamMembers(teamID uuid.UUID, memberIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddContestTeamMembers", teamID, memberIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddContestTeamMembers indicates an expected call of AddContestTeamMembers.
func (mr *MockContestRepositoryMockRecorder) AddContestTeamMembers(teamID, memberIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddContestTeamMembers", reflect.TypeOf((*MockContestRepository)(nil).AddContestTeamMembers), teamID, memberIDs)
}

// CreateContest mocks base method.
func (m *MockContestRepository) CreateContest(args *repository.CreateContestArgs) (*domain.Contest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContest", args)
	ret0, _ := ret[0].(*domain.Contest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContest indicates an expected call of CreateContest.
func (mr *MockContestRepositoryMockRecorder) CreateContest(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContest", reflect.TypeOf((*MockContestRepository)(nil).CreateContest), args)
}

// CreateContestTeam mocks base method.
func (m *MockContestRepository) CreateContestTeam(contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContestTeam", contestID, args)
	ret0, _ := ret[0].(*domain.ContestTeamDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContestTeam indicates an expected call of CreateContestTeam.
func (mr *MockContestRepositoryMockRecorder) CreateContestTeam(contestID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContestTeam", reflect.TypeOf((*MockContestRepository)(nil).CreateContestTeam), contestID, args)
}

// DeleteContest mocks base method.
func (m *MockContestRepository) DeleteContest(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContest", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContest indicates an expected call of DeleteContest.
func (mr *MockContestRepositoryMockRecorder) DeleteContest(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContest", reflect.TypeOf((*MockContestRepository)(nil).DeleteContest), id)
}

// DeleteContestTeamMembers mocks base method.
func (m *MockContestRepository) DeleteContestTeamMembers(teamID uuid.UUID, memberIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContestTeamMembers", teamID, memberIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContestTeamMembers indicates an expected call of DeleteContestTeamMembers.
func (mr *MockContestRepositoryMockRecorder) DeleteContestTeamMembers(teamID, memberIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContestTeamMembers", reflect.TypeOf((*MockContestRepository)(nil).DeleteContestTeamMembers), teamID, memberIDs)
}

// GetContest mocks base method.
func (m *MockContestRepository) GetContest(id uuid.UUID) (*domain.ContestDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContest", id)
	ret0, _ := ret[0].(*domain.ContestDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContest indicates an expected call of GetContest.
func (mr *MockContestRepositoryMockRecorder) GetContest(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContest", reflect.TypeOf((*MockContestRepository)(nil).GetContest), id)
}

// GetContestTeam mocks base method.
func (m *MockContestRepository) GetContestTeam(contestID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContestTeam", contestID, teamID)
	ret0, _ := ret[0].(*domain.ContestTeamDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContestTeam indicates an expected call of GetContestTeam.
func (mr *MockContestRepositoryMockRecorder) GetContestTeam(contestID, teamID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContestTeam", reflect.TypeOf((*MockContestRepository)(nil).GetContestTeam), contestID, teamID)
}

// GetContestTeamMembers mocks base method.
func (m *MockContestRepository) GetContestTeamMembers(contestID, teamID uuid.UUID) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContestTeamMembers", contestID, teamID)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContestTeamMembers indicates an expected call of GetContestTeamMembers.
func (mr *MockContestRepositoryMockRecorder) GetContestTeamMembers(contestID, teamID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContestTeamMembers", reflect.TypeOf((*MockContestRepository)(nil).GetContestTeamMembers), contestID, teamID)
}

// GetContestTeams mocks base method.
func (m *MockContestRepository) GetContestTeams(contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContestTeams", contestID)
	ret0, _ := ret[0].([]*domain.ContestTeam)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContestTeams indicates an expected call of GetContestTeams.
func (mr *MockContestRepositoryMockRecorder) GetContestTeams(contestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContestTeams", reflect.TypeOf((*MockContestRepository)(nil).GetContestTeams), contestID)
}

// GetContests mocks base method.
func (m *MockContestRepository) GetContests() ([]*domain.Contest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContests")
	ret0, _ := ret[0].([]*domain.Contest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContests indicates an expected call of GetContests.
func (mr *MockContestRepositoryMockRecorder) GetContests() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContests", reflect.TypeOf((*MockContestRepository)(nil).GetContests))
}

// UpdateContest mocks base method.
func (m *MockContestRepository) UpdateContest(id uuid.UUID, changes map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContest", id, changes)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContest indicates an expected call of UpdateContest.
func (mr *MockContestRepositoryMockRecorder) UpdateContest(id, changes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContest", reflect.TypeOf((*MockContestRepository)(nil).UpdateContest), id, changes)
}

// UpdateContestTeam mocks base method.
func (m *MockContestRepository) UpdateContestTeam(teamID uuid.UUID, changes map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContestTeam", teamID, changes)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContestTeam indicates an expected call of UpdateContestTeam.
func (mr *MockContestRepositoryMockRecorder) UpdateContestTeam(teamID, changes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContestTeam", reflect.TypeOf((*MockContestRepository)(nil).UpdateContestTeam), teamID, changes)
}
