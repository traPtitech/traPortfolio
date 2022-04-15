package handler

import (
	"errors"
	"fmt"
	"math/rand"
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

// 0 first semester, 1 second semester
func makeYearWithSemester(s int) domain.YearWithSemester {
	return domain.YearWithSemester{
		Year:     random.Time().Year(),
		Semester: s,
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
				sinceSem := rand.Intn(2)
				untilSem := rand.Intn(2)
				duration := domain.YearWithSemesterDuration{
					Since: makeYearWithSemester(sinceSem),
					Until: makeYearWithSemester(untilSem),
				}
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
								Semester: Semester(sinceSem),
								Year:     v.Duration.Since.Year,
							},
							Until: &YearWithSemester{
								Semester: Semester(untilSem),
								Year:     v.Duration.Until.Year,
							},
						},
						Id:   v.ID,
						Name: v.Name,
					})
				}

				s.EXPECT().GetProjects(gomock.Any()).Return(repo, nil)
				return reqBody, "/api/v1/projects"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) ([]*Project, string) {
				s.EXPECT().GetProjects(gomock.Any()).Return(nil, errInternal)
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
				sinceSem := rand.Intn(2)
				untilSem := rand.Intn(2)
				duration := domain.YearWithSemesterDuration{
					Since: makeYearWithSemester(sinceSem),
					Until: makeYearWithSemester(untilSem),
				}
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
						User: User{
							Id:       v.UserID,
							Name:     v.Name,
							RealName: v.RealName,
						},
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
					})
				}
				reqBody := ProjectDetail{
					Project: Project{
						Duration: YearWithSemesterDuration{
							Since: YearWithSemester{
								Semester: Semester(sinceSem),
								Year:     repo.Duration.Since.Year,
							},
							Until: &YearWithSemester{
								Semester: Semester(untilSem),
								Year:     repo.Duration.Until.Year,
							},
						},
						Id:   repo.ID,
						Name: repo.Name,
					},
					Description: repo.Description,
					Link:        repo.Link,
					Members:     members,
				}

				s.EXPECT().GetProject(gomock.Any(), projectID).Return(&repo, nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) (ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(gomock.Any(), projectID).Return(nil, errInternal)
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
