package handler_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

const (
	invalidID = "invalid"
)

func setupContestMock(t *testing.T) (*mock_service.MockContestService, handler.API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockContestService(ctrl)
	api := handler.NewAPI(nil, nil, nil, nil, handler.NewContestHandler(s))

	return s, api
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestContestHandler_GetContests(t *testing.T) {
	until := mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00")

	t.Parallel()
	tests := []struct {
		name         string
		setup        func(s *mock_service.MockContestService, want []*domain.Contest) (path string)
		statusCode   int
		repoContests []*domain.Contest
		hresContests []*handler.Contest
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockContestService, want []*domain.Contest) string {
				s.EXPECT().GetContests(gomock.Any()).Return(want, nil)
				return "/api/v1/contests"
			},
			statusCode: http.StatusOK,
			repoContests: []*domain.Contest{
				{
					ID:        uuid.Nil,
					Name:      "test1",
					TimeStart: mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"),
					TimeEnd:   mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"),
				},
			},
			hresContests: []*handler.Contest{
				{
					Name: "test1",
					Duration: handler.Duration{
						Since: mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"),
						Until: &until,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			for i, v := range tt.hresContests {
				tt.repoContests[i].ID = v.Id
			}
			path := tt.setup(s, tt.repoContests)

			var resBody []*handler.Contest
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, tt.hresContests, resBody)
		})
	}
}

var (
	getContestID = []uuid.UUID{
		uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
		uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
		uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
		uuid.FromStringOrNil("44444444-4444-4444-4444-444444444444"),
		uuid.FromStringOrNil("55555555-5555-5555-5555-555555555555"),
	}
)

func makeContest(t *testing.T) (*domain.ContestDetail, *handler.ContestDetail) {
	t.Helper()

	d := domain.ContestDetail{
		Contest: domain.Contest{
			ID:        getContestID[0],
			Name:      random.AlphaNumericn(rand.Intn(30) + 1),
			TimeStart: random.Time(),
			TimeEnd:   random.Time(),
		},
		Link:        random.RandURLString(),
		Description: random.AlphaNumericn(rand.Intn(30) + 1),
		Teams: []*domain.ContestTeam{
			{
				ID:        getContestID[1],
				ContestID: getContestID[0],
				Name:      random.AlphaNumericn(rand.Intn(30) + 1),
				Result:    random.AlphaNumericn(rand.Intn(30) + 1),
			},
			{
				ID:        getContestID[2],
				ContestID: getContestID[0],
				Name:      random.AlphaNumericn(rand.Intn(30) + 1),
				Result:    random.AlphaNumericn(rand.Intn(30) + 1),
			},
		},
	}

	teams := make([]handler.ContestTeam, len(d.Teams))
	for i, v := range d.Teams {
		teams[i] = handler.ContestTeam{
			Id:     v.ID,
			Name:   v.Name,
			Result: v.Result,
		}
	}

	hres := handler.ContestDetail{
		Contest: handler.Contest{
			Id:   d.ID,
			Name: d.Name,
			Duration: handler.Duration{
				Since: d.TimeStart,
				Until: &d.TimeEnd,
			},
		},
		Link:        d.Link,
		Description: d.Description,
		Teams:       teams,
	}

	return &d, &hres
}

func TestContestHandler_GetContest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (repoContest *domain.ContestDetail, hresContest *handler.ContestDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *handler.ContestDetail, string) {
				want, hres := makeContest(t)
				s.EXPECT().GetContest(gomock.Any(), want.ID).Return(want, nil)
				path := fmt.Sprintf("/api/v1/contests/%s", want.ID.String())

				return want, hres, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *handler.ContestDetail, string) {
				path := fmt.Sprintf("/api/v1/contests/%s", invalidID)
				return &domain.ContestDetail{}, &handler.ContestDetail{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *handler.ContestDetail, string) {
				uid := random.UUID()
				s.EXPECT().GetContest(gomock.Any(), uid).Return(nil, repository.ErrNotFound)

				return &domain.ContestDetail{}, &handler.ContestDetail{}, fmt.Sprintf("/api/v1/contests/%s", uid)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			_, hresContest, path := tt.setup(s)

			var resBody handler.ContestDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresContest, &resBody)
		})
	}
}

func makePostContestRequest(description string, since time.Time, until time.Time, name string, link string) *handler.PostContestJSONRequestBody {
	return &handler.PostContestJSONRequestBody{
		Description: description,
		Duration: handler.Duration{
			Since: since,
			Until: &until,
		},
		Name: name,
		Link: &link,
	}
}

func TestContestHandler_PostContest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string) {
				reqBody = makePostContestRequest(
					random.AlphaNumericn(rand.Intn(30)+1),
					random.Time(),
					random.Time(),
					random.AlphaNumericn(rand.Intn(30)+1),
					random.RandURLString(),
				)
				args := repository.CreateContestArgs{
					Name:        reqBody.Name,
					Description: reqBody.Description,
					Link:        optional.StringFrom(reqBody.Link),
					Since:       reqBody.Duration.Since,
					Until:       optional.TimeFrom(reqBody.Duration.Until),
				}
				want := domain.ContestDetail{
					Contest: domain.Contest{
						ID:        random.UUID(),
						Name:      args.Name,
						TimeStart: args.Since,
						TimeEnd:   args.Until.Time,
					},
					Link:        args.Link.String,
					Description: args.Description,
					Teams:       []*domain.ContestTeam{},
				}
				expectedResBody = &handler.Contest{
					Id:   want.ID,
					Name: want.Name,
					Duration: handler.Duration{
						Since: want.TimeStart,
						Until: &want.TimeEnd,
					},
				}
				s.EXPECT().CreateContest(gomock.Any(), &args).Return(&want, nil)
				path = "/api/v1/contests"
				return reqBody, expectedResBody, &handler.Contest{}, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Bad Request: invalid url",
			setup: func(s *mock_service.MockContestService) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string) {
				reqBody = makePostContestRequest(
					random.AlphaNumericn(rand.Intn(30)+1),
					random.Time(),
					random.Time(),
					random.AlphaNumericn(rand.Intn(30)+1),
					random.AlphaNumericn(rand.Intn(30)+1),
				)
				path = "/api/v1/contests"
				return reqBody, nil, nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Conflict",
			setup: func(s *mock_service.MockContestService) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string) {
				reqBody = makePostContestRequest(
					random.AlphaNumericn(rand.Intn(30)+1),
					random.Time(),
					random.Time(),
					random.AlphaNumericn(rand.Intn(30)+1),
					random.RandURLString(),
				)
				args := repository.CreateContestArgs{
					Name:        reqBody.Name,
					Description: reqBody.Description,
					Link:        optional.StringFrom(reqBody.Link),
					Since:       reqBody.Duration.Since,
					Until:       optional.TimeFrom(reqBody.Duration.Until),
				}
				s.EXPECT().CreateContest(gomock.Any(), &args).Return(nil, repository.ErrAlreadyExists)
				return reqBody, nil, nil, "/api/v1/contests"
			},
			statusCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, res, resBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, res, resBody)
		})
	}
}

func TestContestHandler_PatchContest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *handler.EditContestJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success 1",
			setup: func(s *mock_service.MockContestService) (*handler.EditContestJSONRequestBody, string) {
				contestID := random.UUID()
				name := random.AlphaNumericn(rand.Intn(30) + 1)
				link := random.RandURLString()
				description := random.AlphaNumericn(rand.Intn(30) + 1)
				since := random.Time()
				until := random.Time()
				reqBody := &handler.EditContestJSONRequestBody{
					Name:        &name,
					Link:        &link,
					Description: &description,
					Duration: &handler.Duration{
						Since: since,
						Until: &until,
					},
				}
				args := repository.UpdateContestArgs{
					Name:        optional.StringFrom(reqBody.Name),
					Description: optional.StringFrom(reqBody.Description),
					Link:        optional.StringFrom(reqBody.Link),
					Since:       optional.TimeFrom(&reqBody.Duration.Since),
					Until:       optional.TimeFrom(reqBody.Duration.Until),
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				s.EXPECT().UpdateContest(gomock.Any(), contestID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(s *mock_service.MockContestService) (*handler.EditContestJSONRequestBody, string) {
				path := fmt.Sprintf("/api/v1/contests/%s", invalidID)
				return &handler.EditContestJSONRequestBody{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		// todo validate url

		// {
		// 	name: "BadRequest: Invalid URL",
		// 	setup: func(s *mock_service.MockContestService) (*handler.EditContestJSONRequestBody, string) {
		// 		contestID := random.UUID()
		// 		reqBody := &handler.EditContestJSONRequestBody{
		// 			ContestID:   contestID,
		// 			Name:        ptr(t,random.AlphaNumeric(rand.Intn(30) + 1)),
		// 			Link:        random.AlphaNumeric(rand.Intn(30) + 1),
		// 			Description: random.AlphaNumeric(rand.Intn(30) + 1),
		// 			Duration: handler.OptionalDuration{
		// 				Since: optional.TimeFrom(random.Time()),
		// 				Until: optional.TimeFrom(random.Time()),
		// 			},
		// 		}
		// 		args := repository.UpdateContestArgs{
		// 			Name:        reqBody.Name,
		// 			Description: reqBody.Description,
		// 			Link:        reqBody.Link,
		// 			Since:       reqBody.Duration.Since,
		// 			Until:       reqBody.Duration.Until,
		// 		}
		// 		path := fmt.Sprintf("/api/v1/contests/%s", random.UUID())
		// 		s.EXPECT().UpdateContest(gomock.Any(), contestID, &args).Return(nil)
		// 		return reqBody, path
		// 	},
		// 	statusCode: http.StatusBadRequest,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestContestHandler_DeleteContest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) string {
				contestID := random.UUID()
				s.EXPECT().DeleteContest(gomock.Any(), contestID).Return(nil)
				return fmt.Sprintf("/api/v1/contests/%s", contestID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(s *mock_service.MockContestService) string {
				return fmt.Sprintf("/api/v1/contests/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodDelete, path, nil, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestContestHandler_GetContestTeams(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (hres []*handler.ContestTeam, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (hres []*handler.ContestTeam, path string) {
				contestID := random.UUID()
				repoContestTeams := []*domain.ContestTeam{
					{
						ID:        random.UUID(),
						ContestID: contestID,
						Name:      random.AlphaNumericn(rand.Intn(30) + 1),
						Result:    random.AlphaNumericn(rand.Intn(30) + 1),
					},
					{
						ID:        random.UUID(),
						ContestID: contestID,
						Name:      random.AlphaNumericn(rand.Intn(30) + 1),
						Result:    random.AlphaNumericn(rand.Intn(30) + 1),
					},
				}
				hres = []*handler.ContestTeam{
					{
						Id:     repoContestTeams[0].ID,
						Name:   repoContestTeams[0].Name,
						Result: repoContestTeams[0].Result,
					},
					{
						Id:     repoContestTeams[1].ID,
						Name:   repoContestTeams[1].Name,
						Result: repoContestTeams[1].Result,
					},
				}
				s.EXPECT().GetContestTeams(gomock.Any(), contestID).Return(repoContestTeams, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(s *mock_service.MockContestService) (hres []*handler.ContestTeam, path string) {
				return []*handler.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			expectedHres, path := tt.setup(s)

			hres := make([]*handler.ContestTeam, 0, len(expectedHres))
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestContestHandler_GetContestTeam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (hres handler.ContestTeamDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (handler.ContestTeamDetail, string) {
				teamID := random.UUID()
				contestID := random.UUID()
				repoContestTeamDetail := domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ID:        teamID,
						ContestID: contestID,
						Name:      random.AlphaNumericn(rand.Intn(30) + 1),
						Result:    random.AlphaNumericn(rand.Intn(30) + 1),
					},
					Link:        random.AlphaNumericn(rand.Intn(30) + 1),
					Description: random.AlphaNumericn(rand.Intn(30) + 1),
					Members: []*domain.User{
						{
							ID:       random.UUID(),
							Name:     random.AlphaNumericn(rand.Intn(30) + 1),
							RealName: random.AlphaNumericn(rand.Intn(30) + 1),
						},
						{
							ID:       random.UUID(),
							Name:     random.AlphaNumericn(rand.Intn(30) + 1),
							RealName: random.AlphaNumericn(rand.Intn(30) + 1),
						},
					},
				}
				members := make([]handler.User, 0, len(repoContestTeamDetail.Members))
				for _, member := range repoContestTeamDetail.Members {
					members = append(members, handler.User{
						Id:       member.ID,
						Name:     member.Name,
						RealName: member.RealName,
					})
				}

				hres := handler.ContestTeamDetail{
					ContestTeam: handler.ContestTeam{
						Id:     repoContestTeamDetail.ID,
						Name:   repoContestTeamDetail.Name,
						Result: repoContestTeamDetail.Result,
					},
					Link:        repoContestTeamDetail.Link,
					Description: repoContestTeamDetail.Description,
					Members:     members,
				}

				s.EXPECT().GetContestTeam(gomock.Any(), contestID, teamID).Return(&repoContestTeamDetail, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(s *mock_service.MockContestService) (handler.ContestTeamDetail, string) {
				return handler.ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", invalidID, random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (handler.ContestTeamDetail, string) {
				return handler.ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "NotFound: Contest not found",
			setup: func(s *mock_service.MockContestService) (handler.ContestTeamDetail, string) {
				teamID := random.UUID()
				contestID := random.UUID()
				s.EXPECT().GetContestTeam(gomock.Any(), contestID, teamID).Return(nil, repository.ErrNotFound)
				return handler.ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			expectedHres, path := tt.setup(s)

			var hres handler.ContestTeamDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestContestHandler_PostContestTeam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *handler.PostContestTeamJSONRequestBody, expectedResBody handler.ContestTeam, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*handler.PostContestTeamJSONRequestBody, handler.ContestTeam, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &handler.PostContestTeamJSONRequestBody{
					Name:        random.AlphaNumericn(rand.Intn(30) + 1),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumericn(rand.Intn(30) + 1),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.StringFrom(reqBody.Result),
					Link:        optional.StringFrom(reqBody.Link),
					Description: reqBody.Description,
				}
				want := domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ID:        teamID,
						ContestID: contestID,
						Name:      args.Name,
						Result:    args.Result.String,
					},
					Link:        args.Link.String,
					Description: args.Description,
					Members:     nil,
				}
				expectedResBody := handler.ContestTeam{
					Id:     teamID,
					Name:   want.Name,
					Result: want.Result,
				}
				s.EXPECT().CreateContestTeam(gomock.Any(), contestID, &args).Return(&want, nil)
				return reqBody, expectedResBody, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (*handler.PostContestTeamJSONRequestBody, handler.ContestTeam, string) {
				reqBody := &handler.PostContestTeamJSONRequestBody{
					Name:        random.AlphaNumericn(rand.Intn(30) + 1),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumericn(rand.Intn(30) + 1),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				return reqBody, handler.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) (*handler.PostContestTeamJSONRequestBody, handler.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &handler.PostContestTeamJSONRequestBody{
					Name:        random.AlphaNumericn(rand.Intn(30) + 1),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumericn(rand.Intn(30) + 1),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.StringFrom(reqBody.Result),
					Link:        optional.StringFrom(reqBody.Link),
					Description: reqBody.Description,
				}
				s.EXPECT().CreateContestTeam(gomock.Any(), contestID, &args).Return(nil, repository.ErrNotFound)
				return reqBody, handler.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "conflict contest",
			setup: func(s *mock_service.MockContestService) (*handler.PostContestTeamJSONRequestBody, handler.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &handler.PostContestTeamJSONRequestBody{
					Name:        random.AlphaNumericn(rand.Intn(30) + 1),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumericn(rand.Intn(30) + 1),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.StringFrom(reqBody.Result),
					Link:        optional.StringFrom(reqBody.Link),
					Description: reqBody.Description,
				}
				s.EXPECT().CreateContestTeam(gomock.Any(), contestID, &args).Return(nil, repository.ErrAlreadyExists)
				return reqBody, handler.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, res, path := tt.setup(s)

			var resBody handler.ContestTeam
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, res, resBody)
		})
	}
}

func TestContestHandler_PatchContestTeam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *handler.EditContestTeamJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*handler.EditContestTeamJSONRequestBody, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &handler.EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Description: ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				args := repository.UpdateContestTeamArgs{
					Name:        optional.StringFrom(reqBody.Name),
					Link:        optional.StringFrom(reqBody.Link),
					Result:      optional.StringFrom(reqBody.Result),
					Description: optional.StringFrom(reqBody.Description),
				}
				s.EXPECT().UpdateContestTeam(gomock.Any(), teamID, &args).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (*handler.EditContestTeamJSONRequestBody, string) {
				reqBody := &handler.EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Description: ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", invalidID, random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(s *mock_service.MockContestService) (*handler.EditContestTeamJSONRequestBody, string) {
				reqBody := &handler.EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Description: ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) (*handler.EditContestTeamJSONRequestBody, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &handler.EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
					Description: ptr(t, random.AlphaNumericn(rand.Intn(30)+1)),
				}
				args := repository.UpdateContestTeamArgs{
					Name:        optional.StringFrom(reqBody.Name),
					Link:        optional.StringFrom(reqBody.Link),
					Result:      optional.StringFrom(reqBody.Result),
					Description: optional.StringFrom(reqBody.Description),
				}
				s.EXPECT().UpdateContestTeam(gomock.Any(), teamID, &args).Return(repository.ErrNotFound)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestContestHandler_GetContestTeamMember(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (hres []*handler.User, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) ([]*handler.User, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				users := []*domain.User{
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumericn(rand.Intn(30) + 1),
						RealName: random.AlphaNumericn(rand.Intn(30) + 1),
					},
				}
				hres := make([]*handler.User, len(users))
				for i, user := range users {
					hres[i] = &handler.User{
						Id:       user.ID,
						Name:     user.Name,
						RealName: user.RealName,
					}
				}

				s.EXPECT().GetContestTeamMembers(gomock.Any(), contestID, teamID).Return(users, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) ([]*handler.User, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(s *mock_service.MockContestService) ([]*handler.User, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) ([]*handler.User, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				s.EXPECT().GetContestTeamMembers(gomock.Any(), contestID, teamID).Return(nil, repository.ErrNotFound)
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			expectedHres, path := tt.setup(s)

			var hres []*handler.User
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestContestHandler_PostContestTeamMember(t *testing.T) {
	t.Parallel()

	type Req struct {
		Members []uuid.UUID `json:"members"`
	}
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *Req, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &Req{
					Members: []uuid.UUID{
						random.UUID(),
						random.UUID(),
					},
				}
				s.EXPECT().AddContestTeamMembers(gomock.Any(), teamID, reqBody.Members).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest or team not exist",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &Req{
					Members: []uuid.UUID{
						random.UUID(),
						random.UUID(),
					},
				}
				s.EXPECT().AddContestTeamMembers(gomock.Any(), teamID, reqBody.Members).Return(repository.ErrNotFound)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestContestHandler_DeleteContestTeamMember(t *testing.T) {
	t.Parallel()

	type Req struct {
		Members []uuid.UUID `json:"members"`
	}
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (*Req, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &Req{
					Members: []uuid.UUID{
						random.UUID(),
						random.UUID(),
					},
				}
				s.EXPECT().DeleteContestTeamMembers(gomock.Any(), teamID, reqBody.Members).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest or team not exist",
			setup: func(s *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &Req{
					Members: []uuid.UUID{
						random.UUID(),
						random.UUID(),
					},
				}
				s.EXPECT().DeleteContestTeamMembers(gomock.Any(), teamID, reqBody.Members).Return(repository.ErrNotFound)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodDelete, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
