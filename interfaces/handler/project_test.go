package handler

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/random"
)

var (
	errInternal = errors.New("internal error")
)

func setupProjectMock(t *testing.T) (*mock_service.MockProjectService, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockProjectService(ctrl)
	api := NewAPI(nil, nil, NewProjectHandler(s), nil, nil, nil)

	return s, api
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
		setup      func(s *mock_service.MockProjectService) ([]*schema.Project, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) ([]*schema.Project, string) {
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
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: schema.Semester(v.Duration.Since.Semester),
							},
							Until: &schema.YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: schema.Semester(v.Duration.Until.Semester),
							},
						},
						Id:   v.ID,
						Name: v.Name,
					})
				}

				s.EXPECT().GetProjects(anyCtx{}).Return(repo, nil)
				return reqBody, "/api/v1/projects"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) ([]*schema.Project, string) {
				s.EXPECT().GetProjects(anyCtx{}).Return(nil, errInternal)
				return nil, "/api/v1/projects"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			expectedHres, path := tt.setup(s)

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
		setup      func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
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
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: schema.Semester(v.Duration.Since.Semester),
							},
							Until: &schema.YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: schema.Semester(v.Duration.Until.Semester),
							},
						},
						Id:       v.User.ID,
						Name:     v.User.Name,
						RealName: v.User.RealName(),
					})
				}
				reqBody := &schema.ProjectDetail{
					Description: repo.Description,
					Duration: schema.YearWithSemesterDuration{
						Since: schema.YearWithSemester{
							Semester: schema.Semester(repo.Duration.Since.Semester),
							Year:     repo.Duration.Since.Year,
						},
						Until: &schema.YearWithSemester{
							Semester: schema.Semester(repo.Duration.Until.Semester),
							Year:     repo.Duration.Until.Year,
						},
					},
					Id:      repo.ID,
					Link:    repo.Link,
					Members: members,
					Name:    repo.Name,
				}

				s.EXPECT().GetProject(anyCtx{}, projectID).Return(&repo, nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Bad Request: Validate error: invalid projectID",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				return nil, fmt.Sprintf("/api/v1/projects/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, errInternal)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Validation Error",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				projectID := random.AlphaNumericN(36)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found Error",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, repository.ErrNotFound)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			expectedHres, path := tt.setup(s)

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
		setup      func(s *mock_service.MockProjectService) (reqBody *schema.CreateProjectRequest, expectedResBody schema.Project, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (reqBody *schema.CreateProjectRequest, expectedResBody schema.Project, path string) {
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
						Duration: domain.YearWithSemesterDuration{
							Since: domain.YearWithSemester{
								Year:     args.SinceYear,
								Semester: args.SinceSemester,
							},
							Until: domain.YearWithSemester{
								Year:     args.UntilYear,
								Semester: args.UntilSemester,
							},
						},
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
				s.EXPECT().CreateProject(anyCtx{}, &args).Return(&want, nil)
				return reqBody, expectedResBody, "/api/v1/projects"
			},
			statusCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			reqBody, res, path := tt.setup(s)

			var resBody schema.Project
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, resBody, res)
		})
	}
}
