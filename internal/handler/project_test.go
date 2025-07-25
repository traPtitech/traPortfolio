package handler

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository/mock_repository"
	"go.uber.org/mock/gomock"
)

var (
	errInternal = errors.New("internal error")
)

func setupProjectMock(t *testing.T) (MockRepository, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	project := mock_repository.NewMockProjectRepository(ctrl)
	mr := MockRepository{project: project}
	api := NewAPI(nil, nil, NewProjectHandler(project), nil, nil, nil)

	return mr, api
}

func makeCreateProjectRequest(t *testing.T, description string, since schema.YearWithSemester, until *schema.YearWithSemester, name string, link string) *schema.CreateProjectRequest {
	t.Helper()
	return &schema.CreateProjectRequest{
		Description: description,
		Duration: schema.YearWithSemesterDuration{
			Since: since,
			Until: until,
		},
		Name: name,
		Link: &link,
	}
}

func TestProjectHandler_GetProjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) ([]*schema.Project, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) ([]*schema.Project, string) {
				duration := random.Duration()
				repo := []*domain.Project{
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						Duration: duration,
					},
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						Duration: duration,
					},
				}

				var reqBody []*schema.Project
				for _, v := range repo {
					reqBody = append(reqBody, &schema.Project{
						Duration: schema.ConvertDuration(v.Duration),
						Id:       v.ID,
						Name:     v.Name,
					})
				}

				mr.project.EXPECT().GetProjects(anyCtx{}, &repository.GetProjectsArgs{}).Return(repo, nil)
				return reqBody, "/api/v1/projects"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(mr MockRepository) ([]*schema.Project, string) {
				mr.project.EXPECT().GetProjects(anyCtx{}, &repository.GetProjectsArgs{}).Return(nil, errInternal)
				return nil, "/api/v1/projects"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			r, api := setupProjectMock(t)

			expectedHres, path := tt.setup(r)

			hres := []*schema.Project(nil)
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestProjectHandler_GetProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (*schema.ProjectDetail, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (*schema.ProjectDetail, string) {
				duration := random.Duration()
				projectID := random.UUID()
				repo := domain.ProjectDetail{
					Project: domain.Project{
						ID:       projectID,
						Name:     random.AlphaNumeric(),
						Duration: duration,
					},
					Description: random.AlphaNumeric(),
					Link:        random.RandURLString(),
					Members: []*domain.UserWithDuration{
						{
							User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							Duration: random.Duration(),
						},
					},
				}

				var members []schema.ProjectMember
				for _, v := range repo.Members {
					members = append(members, schema.ProjectMember{
						Duration: schema.ConvertDuration(v.Duration),
						Id:       v.User.ID,
						Name:     v.User.Name,
						RealName: v.User.RealName(),
					})
				}
				reqBody := &schema.ProjectDetail{
					Description: repo.Description,
					Duration:    schema.ConvertDuration(repo.Duration),
					Id:          repo.ID,
					Link:        repo.Link,
					Members:     members,
					Name:        repo.Name,
				}

				mr.project.EXPECT().GetProject(anyCtx{}, projectID).Return(&repo, nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Bad Request: Validate error: invalid projectID",
			setup: func(_ MockRepository) (*schema.ProjectDetail, string) {
				return nil, fmt.Sprintf("/api/v1/projects/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			setup: func(mr MockRepository) (*schema.ProjectDetail, string) {
				projectID := random.UUID()
				mr.project.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, errInternal)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Validation Error",
			setup: func(_ MockRepository) (*schema.ProjectDetail, string) {
				projectID := random.AlphaNumericN(36)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found Error",
			setup: func(mr MockRepository) (*schema.ProjectDetail, string) {
				projectID := random.UUID()
				mr.project.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, repository.ErrNotFound)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			r, api := setupProjectMock(t)

			expectedHres, path := tt.setup(r)

			var hres *schema.ProjectDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestProjectHandler_CreateProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (reqBody *schema.CreateProjectRequest, expectedResBody schema.Project, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (reqBody *schema.CreateProjectRequest, expectedResBody schema.Project, path string) {
				duration := random.Duration()
				reqBody = makeCreateProjectRequest(
					t,
					random.AlphaNumeric(),
					schema.ConvertDuration(duration).Since,
					schema.ConvertDuration(duration).Until,
					random.AlphaNumeric(),
					random.RandURLString(),
				)
				args := repository.CreateProjectArgs{
					Name:          reqBody.Name,
					Description:   reqBody.Description,
					Link:          optional.FromPtr(reqBody.Link),
					SinceYear:     reqBody.Duration.Since.Year,
					SinceSemester: int(reqBody.Duration.Since.Semester),
					UntilYear:     reqBody.Duration.Until.Year,
					UntilSemester: int(reqBody.Duration.Until.Semester),
				}
				want := domain.ProjectDetail{
					Project: domain.Project{
						ID:   random.UUID(),
						Name: args.Name,
						Duration: domain.NewYearWithSemesterDuration(
							args.SinceYear,
							args.SinceSemester,
							args.UntilYear,
							args.UntilSemester,
						),
					},
					Description: args.Description,
					Link:        args.Link.ValueOrZero(),
					Members:     nil,
				}
				expectedResBody = schema.Project{
					Duration: schema.ConvertDuration(want.Duration),
					Id:       want.ID,
					Name:     want.Name,
				}
				mr.project.EXPECT().CreateProject(anyCtx{}, &args).Return(&want, nil)
				return reqBody, expectedResBody, "/api/v1/projects"
			},
			statusCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			r, api := setupProjectMock(t)

			reqBody, res, path := tt.setup(r)

			var resBody schema.Project
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, resBody, res)
		})
	}
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (path string) {
				projectID := random.UUID()
				mr.project.EXPECT().DeleteProject(anyCtx{}, projectID).Return(nil)
				return fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Bad Request: Invalid Project ID",
			setup: func(mr MockRepository) (path string) {
				return fmt.Sprintf("/api/v1/projects/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodDelete, path, nil, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestProjectHandler_EditProjectMembers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (reqBody *schema.EditProjectMembersRequest, path string)
		statusCode int
	}{
		{
			name: "Success: Add Member",
			setup: func(mr MockRepository) (*schema.EditProjectMembersRequest, string) {
				projectID := random.UUID()
				userID := random.UUID()
				userDuration := random.Duration()
				reqBody := &schema.EditProjectMembersRequest{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.ConvertDuration(userDuration),
							UserId:   userID,
						},
					},
				}
				memberReq := []*repository.EditProjectMemberArgs{
					{
						UserID:        userID,
						SinceYear:     userDuration.Since.Year,
						SinceSemester: userDuration.Since.Semester,
						UntilYear:     userDuration.Until.ValueOrZero().Year,
						UntilSemester: userDuration.Until.ValueOrZero().Semester,
					},
				}
				mr.project.EXPECT().EditProjectMembers(anyCtx{}, projectID, memberReq).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Success: Delete All Members",
			setup: func(mr MockRepository) (reqBody *schema.EditProjectMembersRequest, path string) {
				projectID := random.UUID()
				mr.project.EXPECT().EditProjectMembers(anyCtx{}, projectID, []*repository.EditProjectMemberArgs{}).Return(nil)
				return &schema.EditProjectMembersRequest{Members: []schema.MemberIDWithYearWithSemesterDuration{}}, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid Project ID",
			setup: func(mr MockRepository) (reqBody *schema.EditProjectMembersRequest, path string) {
				projectID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: member is empty",
			setup: func(mr MockRepository) (reqBody *schema.EditProjectMembersRequest, path string) {
				projectID := random.UUID()
				return &schema.EditProjectMembersRequest{}, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: memberID is invalid",
			setup: func(mr MockRepository) (reqBody *schema.EditProjectMembersRequest, path string) {
				projectID := random.UUID()
				duration := random.Duration()
				return &schema.EditProjectMembersRequest{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.ConvertDuration(duration),
							UserId:   uuid.Nil,
						},
					},
				}, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: duplicated user",
			setup: func(mr MockRepository) (reqBody *schema.EditProjectMembersRequest, path string) {
				projectID := random.UUID()
				userID := random.UUID()
				duration := random.Duration()
				return &schema.EditProjectMembersRequest{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.YearWithSemesterDuration{
								Since: schema.YearWithSemester{
									Semester: schema.Semester(duration.Since.Semester),
									Year:     duration.Since.Year,
								},
								Until: &schema.YearWithSemester{
									Semester: schema.Semester(duration.Until.ValueOrZero().Semester),
									Year:     duration.Until.ValueOrZero().Year,
								},
							},
							UserId: userID,
						},
						{
							Duration: schema.YearWithSemesterDuration{
								Since: schema.YearWithSemester{
									Semester: schema.Semester(duration.Since.Semester),
									Year:     duration.Since.Year,
								},
								Until: &schema.YearWithSemester{
									Semester: schema.Semester(duration.Until.ValueOrZero().Semester),
									Year:     duration.Until.ValueOrZero().Year,
								},
							},
							UserId: userID,
						},
					},
				}, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: member is already exists",
			setup: func(mr MockRepository) (*schema.EditProjectMembersRequest, string) {
				userID := random.UUID()
				projectID := random.UUID()
				duration := random.Duration()
				reqBody := &schema.EditProjectMembersRequest{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.ConvertDuration(duration),
							UserId:   userID,
						},
					},
				}
				mr.project.EXPECT().EditProjectMembers(anyCtx{}, projectID, []*repository.EditProjectMemberArgs{
					{
						UserID:        userID,
						SinceYear:     int(duration.Since.Year),
						SinceSemester: int(duration.Since.Semester),
						UntilYear:     int(duration.Until.ValueOrZero().Year),
						UntilSemester: int(duration.Until.ValueOrZero().Semester),
					},
				}).Return(repository.ErrInvalidArg)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPut, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
