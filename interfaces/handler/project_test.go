package handler_test

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/random"
)

var (
	errInternal = errors.New("internal error")
)

// 0 first semester, 1 second semester
func makeYearWithSemester(s int) domain.YearWithSemester {
	return domain.YearWithSemester{
		Year:     random.Time().Year(),
		Semester: s,
	}
}

func TestProjecttHandler_GetAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) ([]*handler.Project, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(th *handler.TestHandlers) ([]*handler.Project, string) {
				sinceSem := rand.Intn(2)
				untilSem := rand.Intn(2)
				duration := domain.YearWithSemesterDuration{
					Since: makeYearWithSemester(sinceSem),
					Until: makeYearWithSemester(untilSem),
				}
				repo := []*domain.Project{
					{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(rand.Intn(30) + 1),
						Duration:    duration,
						Description: random.AlphaNumeric(rand.Intn(30) + 1),
						Link:        random.RandURLString(),
						Members: []*domain.ProjectMember{
							{
								UserID:   random.UUID(),
								Name:     random.AlphaNumeric(rand.Intn(30) + 1),
								RealName: random.AlphaNumeric(rand.Intn(30) + 1),
								Duration: random.Duration(),
							},
						},
					},
					{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(rand.Intn(30) + 1),
						Duration:    duration,
						Description: random.AlphaNumeric(rand.Intn(30) + 1),
						Link:        random.RandURLString(),
						Members: []*domain.ProjectMember{
							{
								UserID:   random.UUID(),
								Name:     random.AlphaNumeric(rand.Intn(30) + 1),
								RealName: random.AlphaNumeric(rand.Intn(30) + 1),
								Duration: random.Duration(),
							},
						},
					},
				}

				var reqBody []*handler.Project
				for _, v := range repo {
					reqBody = append(reqBody, &handler.Project{
						Duration: handler.YearWithSemesterDuration{
							Since: handler.YearWithSemester{
								Semester: handler.Semester(sinceSem),
								Year:     v.Duration.Since.Year,
							},
							Until: &handler.YearWithSemester{
								Semester: handler.Semester(untilSem),
								Year:     v.Duration.Until.Year,
							},
						},
						Id:   v.ID,
						Name: v.Name,
					})
				}

				th.Service.MockProjectService.EXPECT().GetProjects(gomock.Any()).Return(repo, nil)
				return reqBody, "/api/v1/projects"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(th *handler.TestHandlers) ([]*handler.Project, string) {
				th.Service.MockProjectService.EXPECT().GetProjects(gomock.Any()).Return(nil, errInternal)
				return nil, "/api/v1/projects"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			expectedHres, path := tt.setup(&handlers)

			hres := []*handler.Project(nil)
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestProjecttHandler_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (handler.ProjectDetail, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(th *handler.TestHandlers) (handler.ProjectDetail, string) {
				sinceSem := rand.Intn(2)
				untilSem := rand.Intn(2)
				duration := domain.YearWithSemesterDuration{
					Since: makeYearWithSemester(sinceSem),
					Until: makeYearWithSemester(untilSem),
				}
				projectID := random.UUID()
				repo := domain.Project{
					ID:          projectID,
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Duration:    duration,
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Members: []*domain.ProjectMember{
						{
							UserID:   random.UUID(),
							Name:     random.AlphaNumeric(rand.Intn(30) + 1),
							RealName: random.AlphaNumeric(rand.Intn(30) + 1),
							Duration: random.Duration(),
						},
					},
				}

				var members []handler.ProjectMember
				for _, v := range repo.Members {
					members = append(members, handler.ProjectMember{
						User: handler.User{
							Id:       v.UserID,
							Name:     v.Name,
							RealName: v.RealName,
						},
						Duration: handler.YearWithSemesterDuration{
							Since: handler.YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: handler.Semester(v.Duration.Since.Semester),
							},
							Until: &handler.YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: handler.Semester(v.Duration.Until.Semester),
							},
						},
					})
				}
				reqBody := handler.ProjectDetail{
					Project: handler.Project{
						Duration: handler.YearWithSemesterDuration{
							Since: handler.YearWithSemester{
								Semester: handler.Semester(sinceSem),
								Year:     repo.Duration.Since.Year,
							},
							Until: &handler.YearWithSemester{
								Semester: handler.Semester(untilSem),
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

				th.Service.MockProjectService.EXPECT().GetProject(gomock.Any(), projectID).Return(&repo, nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(th *handler.TestHandlers) (handler.ProjectDetail, string) {
				projectID := random.UUID()
				th.Service.MockProjectService.EXPECT().GetProject(gomock.Any(), projectID).Return(nil, errInternal)
				return handler.ProjectDetail{}, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			expectedHres, path := tt.setup(&handlers)

			hres := handler.ProjectDetail{}
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}
