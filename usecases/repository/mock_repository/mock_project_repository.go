// Code generated by MockGen. DO NOT EDIT.
// Source: project_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/traPortfolio/domain"
	repository "github.com/traPtitech/traPortfolio/usecases/repository"
)

// MockProjectRepository is a mock of ProjectRepository interface.
type MockProjectRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProjectRepositoryMockRecorder
}

// MockProjectRepositoryMockRecorder is the mock recorder for MockProjectRepository.
type MockProjectRepositoryMockRecorder struct {
	mock *MockProjectRepository
}

// NewMockProjectRepository creates a new mock instance.
func NewMockProjectRepository(ctrl *gomock.Controller) *MockProjectRepository {
	mock := &MockProjectRepository{ctrl: ctrl}
	mock.recorder = &MockProjectRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectRepository) EXPECT() *MockProjectRepositoryMockRecorder {
	return m.recorder
}

// AddProjectMembers mocks base method.
func (m *MockProjectRepository) AddProjectMembers(projectID uuid.UUID, args []*repository.CreateProjectMemberArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProjectMembers", projectID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProjectMembers indicates an expected call of AddProjectMembers.
func (mr *MockProjectRepositoryMockRecorder) AddProjectMembers(projectID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProjectMembers", reflect.TypeOf((*MockProjectRepository)(nil).AddProjectMembers), projectID, args)
}

// CreateProject mocks base method.
func (m *MockProjectRepository) CreateProject(args *repository.CreateProjectArgs) (*domain.ProjectDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", args)
	ret0, _ := ret[0].(*domain.ProjectDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectRepositoryMockRecorder) CreateProject(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectRepository)(nil).CreateProject), args)
}

// DeleteProjectMembers mocks base method.
func (m *MockProjectRepository) DeleteProjectMembers(projectID uuid.UUID, memberIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProjectMembers", projectID, memberIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProjectMembers indicates an expected call of DeleteProjectMembers.
func (mr *MockProjectRepositoryMockRecorder) DeleteProjectMembers(projectID, memberIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProjectMembers", reflect.TypeOf((*MockProjectRepository)(nil).DeleteProjectMembers), projectID, memberIDs)
}

// GetProject mocks base method.
func (m *MockProjectRepository) GetProject(projectID uuid.UUID) (*domain.ProjectDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", projectID)
	ret0, _ := ret[0].(*domain.ProjectDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject.
func (mr *MockProjectRepositoryMockRecorder) GetProject(projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectRepository)(nil).GetProject), projectID)
}

// GetProjectMembers mocks base method.
func (m *MockProjectRepository) GetProjectMembers(projectID uuid.UUID) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectMembers", projectID)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectMembers indicates an expected call of GetProjectMembers.
func (mr *MockProjectRepositoryMockRecorder) GetProjectMembers(projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectMembers", reflect.TypeOf((*MockProjectRepository)(nil).GetProjectMembers), projectID)
}

// GetProjects mocks base method.
func (m *MockProjectRepository) GetProjects() ([]*domain.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjects")
	ret0, _ := ret[0].([]*domain.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjects indicates an expected call of GetProjects.
func (mr *MockProjectRepositoryMockRecorder) GetProjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjects", reflect.TypeOf((*MockProjectRepository)(nil).GetProjects))
}

// UpdateProject mocks base method.
func (m *MockProjectRepository) UpdateProject(projectID uuid.UUID, args *repository.UpdateProjectArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProject", projectID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProject indicates an expected call of UpdateProject.
func (mr *MockProjectRepositoryMockRecorder) UpdateProject(projectID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockProjectRepository)(nil).UpdateProject), projectID, args)
}
