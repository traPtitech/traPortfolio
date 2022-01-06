package handler_test

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

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
func makeSemesterTime(s int) time.Time {
	t := random.Time()
	var m time.Month
	if s == 0 {
		m = time.August
	} else {
		m = time.December
	}
	newT := time.Date(t.Year(), m, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return newT
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
				repo := []*domain.Project{
					{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(rand.Intn(30) + 1),
						Since:       makeSemesterTime(sinceSem),
						Until:       makeSemesterTime(untilSem),
						Description: random.AlphaNumeric(rand.Intn(30) + 1),
						Link:        random.RandURLString(),
						Members: []*domain.ProjectMember{
							{
								UserID:   random.UUID(),
								Name:     random.AlphaNumeric(rand.Intn(30) + 1),
								RealName: random.AlphaNumeric(rand.Intn(30) + 1),
								Since:    random.Time(),
								Until:    random.Time(),
							},
						},
					},
					{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(rand.Intn(30) + 1),
						Since:       makeSemesterTime(sinceSem),
						Until:       makeSemesterTime(untilSem),
						Description: random.AlphaNumeric(rand.Intn(30) + 1),
						Link:        random.RandURLString(),
						Members: []*domain.ProjectMember{
							{
								UserID:   random.UUID(),
								Name:     random.AlphaNumeric(rand.Intn(30) + 1),
								RealName: random.AlphaNumeric(rand.Intn(30) + 1),
								Since:    random.Time(),
								Until:    random.Time(),
							},
						},
					},
				}

				var reqBody []*handler.Project
				for _, v := range repo {
					reqBody = append(reqBody, &handler.Project{
						Duration: handler.ProjectDuration{
							Since: handler.YearWithSemester{
								Semester: handler.Semester(sinceSem),
								Year:     v.Since.Year(),
							},
							Until: &handler.YearWithSemester{
								Semester: handler.Semester(untilSem),
								Year:     v.Until.Year(),
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
