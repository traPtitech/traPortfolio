// Code generated by MockGen. DO NOT EDIT.
// Source: project.go

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

// MockProjectService is a mock of ProjectService interface.
type MockProjectService struct {
	ctrl     *gomock.Controller
	recorder *MockProjectServiceMockRecorder
}

// MockProjectServiceMockRecorder is the mock recorder for MockProjectService.
type MockProjectServiceMockRecorder struct {
	mock *MockProjectService
}

// NewMockProjectService creates a new mock instance.
func NewMockProjectService(ctrl *gomock.Controller) *MockProjectService {
	mock := &MockProjectService{ctrl: ctrl}
	mock.recorder = &MockProjectServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectService) EXPECT() *MockProjectServiceMockRecorder {
	return m.recorder
}

// AddProjectMembers mocks base method.
func (m *MockProjectService) AddProjectMembers(ctx context.Context, projectID uuid.UUID, args []*repository.CreateProjectMemberArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProjectMembers", ctx, projectID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProjectMembers indicates an expected call of AddProjectMembers.
func (mr *MockProjectServiceMockRecorder) AddProjectMembers(ctx, projectID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProjectMembers", reflect.TypeOf((*MockProjectService)(nil).AddProjectMembers), ctx, projectID, args)
}

// CreateProject mocks base method.
func (m *MockProjectService) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.ProjectDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", ctx, args)
	ret0, _ := ret[0].(*domain.ProjectDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectServiceMockRecorder) CreateProject(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectService)(nil).CreateProject), ctx, args)
}

// DeleteProjectMembers mocks base method.
func (m *MockProjectService) DeleteProjectMembers(ctx context.Context, projectID uuid.UUID, memberIDs []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProjectMembers", ctx, projectID, memberIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProjectMembers indicates an expected call of DeleteProjectMembers.
func (mr *MockProjectServiceMockRecorder) DeleteProjectMembers(ctx, projectID, memberIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProjectMembers", reflect.TypeOf((*MockProjectService)(nil).DeleteProjectMembers), ctx, projectID, memberIDs)
}

// GetProject mocks base method.
func (m *MockProjectService) GetProject(ctx context.Context, projectID uuid.UUID) (*domain.ProjectDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", ctx, projectID)
	ret0, _ := ret[0].(*domain.ProjectDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject.
func (mr *MockProjectServiceMockRecorder) GetProject(ctx, projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectService)(nil).GetProject), ctx, projectID)
}

// GetProjectMembers mocks base method.
func (m *MockProjectService) GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*domain.UserWithDuration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectMembers", ctx, projectID)
	ret0, _ := ret[0].([]*domain.UserWithDuration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectMembers indicates an expected call of GetProjectMembers.
func (mr *MockProjectServiceMockRecorder) GetProjectMembers(ctx, projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectMembers", reflect.TypeOf((*MockProjectService)(nil).GetProjectMembers), ctx, projectID)
}

// GetProjects mocks base method.
func (m *MockProjectService) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjects", ctx)
	ret0, _ := ret[0].([]*domain.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjects indicates an expected call of GetProjects.
func (mr *MockProjectServiceMockRecorder) GetProjects(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjects", reflect.TypeOf((*MockProjectService)(nil).GetProjects), ctx)
}

// UpdateProject mocks base method.
func (m *MockProjectService) UpdateProject(ctx context.Context, projectID uuid.UUID, args *repository.UpdateProjectArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProject", ctx, projectID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProject indicates an expected call of UpdateProject.
func (mr *MockProjectServiceMockRecorder) UpdateProject(ctx, projectID, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockProjectService)(nil).UpdateProject), ctx, projectID, args)
}
