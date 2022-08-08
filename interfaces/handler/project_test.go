package handler

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
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

func TestProjectHandler_GetProjects(t *testing.T) {
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

func TestProjectHandler_GetProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) (*ProjectDetail, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (*ProjectDetail, string) {
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
				reqBody := &ProjectDetail{
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
			name: "Bad Request: Validate error: invalid projectID",
			setup: func(s *mock_service.MockProjectService) (*ProjectDetail, string) {
				return nil, fmt.Sprintf("/api/v1/projects/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) (*ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, errInternal)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Validation Error",
			setup: func(s *mock_service.MockProjectService) (*ProjectDetail, string) {
				projectID := random.AlphaNumericn(36)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found Error",
			setup: func(s *mock_service.MockProjectService) (*ProjectDetail, string) {
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

			var hres *ProjectDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}
