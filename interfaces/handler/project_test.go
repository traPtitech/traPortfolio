package handler

import (
	"errors"
	"fmt"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"net/http"
	"testing"

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

func makeCreateProjectRequest(description string, since YearWithSemester, until YearWithSemester, name string, link string) *CreateProjectJSONRequestBody {
	return &CreateProjectJSONRequestBody{
		Description: description,
		Duration: YearWithSemesterDuration{
			Since: since,
			Until: &until,
		},
		Name: name,
		Link: &link,
	}
}

func TestProjectHandler_GetAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) ([]*Project, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) ([]*Project, string) {
				duration := random.Duration()
				repo := []*domain.Project{
					{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(),
						Duration:    duration,
						Description: random.AlphaNumeric(),
						Link:        random.RandURLString(),
						Members: []*domain.ProjectMember{
							{
								UserID:   random.UUID(),
								Name:     random.AlphaNumeric(),
								RealName: random.AlphaNumeric(),
								Duration: random.Duration(),
							},
						},
					},
					{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(),
						Duration:    duration,
						Description: random.AlphaNumeric(),
						Link:        random.RandURLString(),
						Members: []*domain.ProjectMember{
							{
								UserID:   random.UUID(),
								Name:     random.AlphaNumeric(),
								RealName: random.AlphaNumeric(),
								Duration: random.Duration(),
							},
						},
					},
				}

				var reqBody []*Project
				for _, v := range repo {
					reqBody = append(reqBody, &Project{
						Duration: YearWithSemesterDuration{
							Since: YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: Semester(v.Duration.Since.Semester),
							},
							Until: &YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: Semester(v.Duration.Until.Semester),
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
			setup: func(s *mock_service.MockProjectService) ([]*Project, string) {
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

			hres := []*Project(nil)
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestProjectHandler_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) (ProjectDetail, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (ProjectDetail, string) {
				duration := random.Duration()
				projectID := random.UUID()
				repo := domain.Project{
					ID:          projectID,
					Name:        random.AlphaNumeric(),
					Duration:    duration,
					Description: random.AlphaNumeric(),
					Link:        random.RandURLString(),
					Members: []*domain.ProjectMember{
						{
							UserID:   random.UUID(),
							Name:     random.AlphaNumeric(),
							RealName: random.AlphaNumeric(),
							Duration: random.Duration(),
						},
					},
				}

				var members []ProjectMember
				for _, v := range repo.Members {
					members = append(members, ProjectMember{
						Duration: YearWithSemesterDuration{
							Since: YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: Semester(v.Duration.Since.Semester),
							},
							Until: &YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: Semester(v.Duration.Until.Semester),
							},
						},
						Id:       v.UserID,
						Name:     v.Name,
						RealName: v.RealName,
					})
				}
				reqBody := ProjectDetail{
					Description: repo.Description,
					Duration: YearWithSemesterDuration{
						Since: YearWithSemester{
							Semester: Semester(repo.Duration.Since.Semester),
							Year:     repo.Duration.Since.Year,
						},
						Until: &YearWithSemester{
							Semester: Semester(repo.Duration.Until.Semester),
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
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) (ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, errInternal)
				return ProjectDetail{}, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			expectedHres, path := tt.setup(s)

			hres := ProjectDetail{}
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
		setup      func(s *mock_service.MockProjectService) (reqBody *CreateProjectJSONRequestBody, expectedResBody Project, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (reqBody *CreateProjectJSONRequestBody, expectedResBody Project, path string) {
				duration := random.Duration()
				reqBody = makeCreateProjectRequest(
					random.AlphaNumeric(),
					convertDuration(duration).Since,
					*convertDuration(duration).Until,
					random.AlphaNumeric(),
					random.RandURLString(),
				)
				args := repository.CreateProjectArgs{
					Name:          reqBody.Name,
					Description:   reqBody.Description,
					Link:          optional.StringFrom(reqBody.Link),
					SinceYear:     reqBody.Duration.Since.Year,
					SinceSemester: int(reqBody.Duration.Since.Semester),
					UntilYear:     reqBody.Duration.Until.Year,
					UntilSemester: int(reqBody.Duration.Until.Semester),
				}
				want := domain.Project{
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
					Description: args.Description,
					Link:        args.Link.String,
					Members:     nil,
				}
				expectedResBody = Project{
					Duration: convertDuration(want.Duration),
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

			var resBody Project
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, resBody, res)
		})
	}
}
