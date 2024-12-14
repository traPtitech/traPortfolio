// Code generated by MockGen. DO NOT EDIT.
// Source: project_repository.go
//
// Generated by this command:
//
//	mockgen -typed -source=project_repository.go -destination=mock_repository/mock_project_repository.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	domain "github.com/traPtitech/traPortfolio/internal/domain"
	repository "github.com/traPtitech/traPortfolio/internal/usecases/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockProjectRepository is a mock of ProjectRepository interface.
type MockProjectRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProjectRepositoryMockRecorder
	isgomock struct{}
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

// CreateProject mocks base method.
func (m *MockProjectRepository) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.ProjectDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", ctx, args)
	ret0, _ := ret[0].(*domain.ProjectDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectRepositoryMockRecorder) CreateProject(ctx, args any) *MockProjectRepositoryCreateProjectCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectRepository)(nil).CreateProject), ctx, args)
	return &MockProjectRepositoryCreateProjectCall{Call: call}
}

// MockProjectRepositoryCreateProjectCall wrap *gomock.Call
type MockProjectRepositoryCreateProjectCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryCreateProjectCall) Return(arg0 *domain.ProjectDetail, arg1 error) *MockProjectRepositoryCreateProjectCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryCreateProjectCall) Do(f func(context.Context, *repository.CreateProjectArgs) (*domain.ProjectDetail, error)) *MockProjectRepositoryCreateProjectCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryCreateProjectCall) DoAndReturn(f func(context.Context, *repository.CreateProjectArgs) (*domain.ProjectDetail, error)) *MockProjectRepositoryCreateProjectCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteProject mocks base method.
func (m *MockProjectRepository) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", ctx, projectID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject.
func (mr *MockProjectRepositoryMockRecorder) DeleteProject(ctx, projectID any) *MockProjectRepositoryDeleteProjectCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProjectRepository)(nil).DeleteProject), ctx, projectID)
	return &MockProjectRepositoryDeleteProjectCall{Call: call}
}

// MockProjectRepositoryDeleteProjectCall wrap *gomock.Call
type MockProjectRepositoryDeleteProjectCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryDeleteProjectCall) Return(arg0 error) *MockProjectRepositoryDeleteProjectCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryDeleteProjectCall) Do(f func(context.Context, uuid.UUID) error) *MockProjectRepositoryDeleteProjectCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryDeleteProjectCall) DoAndReturn(f func(context.Context, uuid.UUID) error) *MockProjectRepositoryDeleteProjectCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// EditProjectMembers mocks base method.
func (m *MockProjectRepository) EditProjectMembers(ctx context.Context, projectID uuid.UUID, args []*repository.EditProjectMemberArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProjectMembers", ctx, projectID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditProjectMembers indicates an expected call of EditProjectMembers.
func (mr *MockProjectRepositoryMockRecorder) EditProjectMembers(ctx, projectID, args any) *MockProjectRepositoryEditProjectMembersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProjectMembers", reflect.TypeOf((*MockProjectRepository)(nil).EditProjectMembers), ctx, projectID, args)
	return &MockProjectRepositoryEditProjectMembersCall{Call: call}
}

// MockProjectRepositoryEditProjectMembersCall wrap *gomock.Call
type MockProjectRepositoryEditProjectMembersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryEditProjectMembersCall) Return(arg0 error) *MockProjectRepositoryEditProjectMembersCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryEditProjectMembersCall) Do(f func(context.Context, uuid.UUID, []*repository.EditProjectMemberArgs) error) *MockProjectRepositoryEditProjectMembersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryEditProjectMembersCall) DoAndReturn(f func(context.Context, uuid.UUID, []*repository.EditProjectMemberArgs) error) *MockProjectRepositoryEditProjectMembersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetProject mocks base method.
func (m *MockProjectRepository) GetProject(ctx context.Context, projectID uuid.UUID) (*domain.ProjectDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", ctx, projectID)
	ret0, _ := ret[0].(*domain.ProjectDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject.
func (mr *MockProjectRepositoryMockRecorder) GetProject(ctx, projectID any) *MockProjectRepositoryGetProjectCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectRepository)(nil).GetProject), ctx, projectID)
	return &MockProjectRepositoryGetProjectCall{Call: call}
}

// MockProjectRepositoryGetProjectCall wrap *gomock.Call
type MockProjectRepositoryGetProjectCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryGetProjectCall) Return(arg0 *domain.ProjectDetail, arg1 error) *MockProjectRepositoryGetProjectCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryGetProjectCall) Do(f func(context.Context, uuid.UUID) (*domain.ProjectDetail, error)) *MockProjectRepositoryGetProjectCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryGetProjectCall) DoAndReturn(f func(context.Context, uuid.UUID) (*domain.ProjectDetail, error)) *MockProjectRepositoryGetProjectCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetProjectMembers mocks base method.
func (m *MockProjectRepository) GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*domain.UserWithDuration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectMembers", ctx, projectID)
	ret0, _ := ret[0].([]*domain.UserWithDuration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectMembers indicates an expected call of GetProjectMembers.
func (mr *MockProjectRepositoryMockRecorder) GetProjectMembers(ctx, projectID any) *MockProjectRepositoryGetProjectMembersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectMembers", reflect.TypeOf((*MockProjectRepository)(nil).GetProjectMembers), ctx, projectID)
	return &MockProjectRepositoryGetProjectMembersCall{Call: call}
}

// MockProjectRepositoryGetProjectMembersCall wrap *gomock.Call
type MockProjectRepositoryGetProjectMembersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryGetProjectMembersCall) Return(arg0 []*domain.UserWithDuration, arg1 error) *MockProjectRepositoryGetProjectMembersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryGetProjectMembersCall) Do(f func(context.Context, uuid.UUID) ([]*domain.UserWithDuration, error)) *MockProjectRepositoryGetProjectMembersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryGetProjectMembersCall) DoAndReturn(f func(context.Context, uuid.UUID) ([]*domain.UserWithDuration, error)) *MockProjectRepositoryGetProjectMembersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetProjects mocks base method.
func (m *MockProjectRepository) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjects", ctx)
	ret0, _ := ret[0].([]*domain.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjects indicates an expected call of GetProjects.
func (mr *MockProjectRepositoryMockRecorder) GetProjects(ctx any) *MockProjectRepositoryGetProjectsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjects", reflect.TypeOf((*MockProjectRepository)(nil).GetProjects), ctx)
	return &MockProjectRepositoryGetProjectsCall{Call: call}
}

// MockProjectRepositoryGetProjectsCall wrap *gomock.Call
type MockProjectRepositoryGetProjectsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryGetProjectsCall) Return(arg0 []*domain.Project, arg1 error) *MockProjectRepositoryGetProjectsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryGetProjectsCall) Do(f func(context.Context) ([]*domain.Project, error)) *MockProjectRepositoryGetProjectsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryGetProjectsCall) DoAndReturn(f func(context.Context) ([]*domain.Project, error)) *MockProjectRepositoryGetProjectsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateProject mocks base method.
func (m *MockProjectRepository) UpdateProject(ctx context.Context, projectID uuid.UUID, args *repository.UpdateProjectArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProject", ctx, projectID, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProject indicates an expected call of UpdateProject.
func (mr *MockProjectRepositoryMockRecorder) UpdateProject(ctx, projectID, args any) *MockProjectRepositoryUpdateProjectCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockProjectRepository)(nil).UpdateProject), ctx, projectID, args)
	return &MockProjectRepositoryUpdateProjectCall{Call: call}
}

// MockProjectRepositoryUpdateProjectCall wrap *gomock.Call
type MockProjectRepositoryUpdateProjectCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProjectRepositoryUpdateProjectCall) Return(arg0 error) *MockProjectRepositoryUpdateProjectCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProjectRepositoryUpdateProjectCall) Do(f func(context.Context, uuid.UUID, *repository.UpdateProjectArgs) error) *MockProjectRepositoryUpdateProjectCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProjectRepositoryUpdateProjectCall) DoAndReturn(f func(context.Context, uuid.UUID, *repository.UpdateProjectArgs) error) *MockProjectRepositoryUpdateProjectCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
