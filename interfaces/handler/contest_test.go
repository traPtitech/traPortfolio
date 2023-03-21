package handler

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

const (
	invalidID = "invalid"
)

func setupContestMock(t *testing.T) (*mock_service.MockContestService, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockContestService(ctrl)
	api := NewAPI(nil, nil, nil, nil, NewContestHandler(s), nil)

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
		hresContests []*Contest
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockContestService, want []*domain.Contest) string {
				s.EXPECT().GetContests(anyCtx{}).Return(want, nil)
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
			hresContests: []*Contest{
				{
					Name: "test1",
					Duration: Duration{
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

			var resBody []*Contest
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, tt.hresContests, resBody)
		})
	}
}

var (
	getContestID = []uuid.UUID{
		uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111"),
		uuid.FromStringOrNil("22222222-2222-4222-8222-222222222222"),
		uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333"),
		uuid.FromStringOrNil("44444444-4444-4444-8444-444444444444"),
		uuid.FromStringOrNil("55555555-5555-4555-8555-555555555555"),
	}
)

func makeContest(t *testing.T) (*domain.ContestDetail, *ContestDetail) {
	t.Helper()

	since, until := random.SinceAndUntil()
	d := domain.ContestDetail{
		Contest: domain.Contest{
			ID:        getContestID[0],
			Name:      random.AlphaNumeric(),
			TimeStart: since,
			TimeEnd:   until,
		},
		Link:        random.RandURLString(),
		Description: random.AlphaNumeric(),
		ContestTeams: []*domain.ContestTeam{
			{
				ID:        getContestID[1],
				ContestID: getContestID[0],
				Name:      random.AlphaNumeric(),
				Result:    random.AlphaNumeric(),
				Members:   make([]*domain.User, 0),
			},
			{
				ID:        getContestID[2],
				ContestID: getContestID[0],
				Name:      random.AlphaNumeric(),
				Result:    random.AlphaNumeric(),
				Members:   make([]*domain.User, 0),
			},
		},
	}

	teams := make([]ContestTeam, len(d.ContestTeams))
	for i, v := range d.ContestTeams {
		teams[i] = ContestTeam{
			Id:      v.ID,
			Members: make([]User, 0),
			Name:    v.Name,
			Result:  v.Result,
		}
	}

	hres := ContestDetail{
		Description: d.Description,
		Duration: Duration{
			Since: d.TimeStart,
			Until: &d.TimeEnd,
		},
		Id:    d.ID,
		Link:  d.Link,
		Name:  d.Name,
		Teams: teams,
	}

	return &d, &hres
}

func TestContestHandler_GetContest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (repoContest *domain.ContestDetail, hresContest *ContestDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *ContestDetail, string) {
				want, hres := makeContest(t)
				s.EXPECT().GetContest(anyCtx{}, want.ID).Return(want, nil)
				path := fmt.Sprintf("/api/v1/contests/%s", want.ID.String())

				return want, hres, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			setup: func(_ *mock_service.MockContestService) (*domain.ContestDetail, *ContestDetail, string) {
				path := fmt.Sprintf("/api/v1/contests/%s", invalidID)
				return &domain.ContestDetail{}, &ContestDetail{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *ContestDetail, string) {
				uid := random.UUID()
				s.EXPECT().GetContest(anyCtx{}, uid).Return(nil, repository.ErrNotFound)

				return &domain.ContestDetail{}, &ContestDetail{}, fmt.Sprintf("/api/v1/contests/%s", uid)
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

			var resBody ContestDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresContest, &resBody)
		})
	}
}

func makeCreateContestRequest(description string, since time.Time, until time.Time, name string, link string) *CreateContestJSONRequestBody {
	return &CreateContestJSONRequestBody{
		Description: description,
		Duration: Duration{
			Since: since,
			Until: &until,
		},
		Name: name,
		Link: &link,
	}
}

func TestContestHandler_CreateContest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *CreateContestJSONRequestBody, expectedResBody *Contest, resBody *Contest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (reqBody *CreateContestJSONRequestBody, expectedResBody *Contest, resBody *Contest, path string) {
				since, until := random.SinceAndUntil()
				reqBody = makeCreateContestRequest(
					random.AlphaNumeric(),
					since,
					until,
					random.AlphaNumeric(),
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
					Link:         args.Link.String,
					Description:  args.Description,
					ContestTeams: []*domain.ContestTeam{},
				}
				expectedResBody = &Contest{
					Id:   want.ID,
					Name: want.Name,
					Duration: Duration{
						Since: want.TimeStart,
						Until: &want.TimeEnd,
					},
				}
				s.EXPECT().CreateContest(anyCtx{}, &args).Return(&want, nil)
				path = "/api/v1/contests"
				return reqBody, expectedResBody, &Contest{}, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Bad Request: invalid url",
			setup: func(_ *mock_service.MockContestService) (reqBody *CreateContestJSONRequestBody, expectedResBody *Contest, resBody *Contest, path string) {
				since, until := random.SinceAndUntil()
				reqBody = makeCreateContestRequest(
					random.AlphaNumeric(),
					since,
					until,
					random.AlphaNumeric(),
					random.AlphaNumeric(),
				)
				path = "/api/v1/contests"
				return reqBody, nil, nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Conflict",
			setup: func(s *mock_service.MockContestService) (reqBody *CreateContestJSONRequestBody, expectedResBody *Contest, resBody *Contest, path string) {
				since, until := random.SinceAndUntil()
				reqBody = makeCreateContestRequest(
					random.AlphaNumeric(),
					since,
					until,
					random.AlphaNumeric(),
					random.RandURLString(),
				)
				args := repository.CreateContestArgs{
					Name:        reqBody.Name,
					Description: reqBody.Description,
					Link:        optional.StringFrom(reqBody.Link),
					Since:       reqBody.Duration.Since,
					Until:       optional.TimeFrom(reqBody.Duration.Until),
				}
				s.EXPECT().CreateContest(anyCtx{}, &args).Return(nil, repository.ErrAlreadyExists)
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
		setup      func(s *mock_service.MockContestService) (reqBody *EditContestJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success 1",
			setup: func(s *mock_service.MockContestService) (*EditContestJSONRequestBody, string) {
				contestID := random.UUID()
				name := random.AlphaNumeric()
				link := random.RandURLString()
				description := random.AlphaNumeric()
				since, until := random.SinceAndUntil()
				reqBody := &EditContestJSONRequestBody{
					Name:        &name,
					Link:        &link,
					Description: &description,
					Duration: &Duration{
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
				s.EXPECT().UpdateContest(anyCtx{}, contestID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(_ *mock_service.MockContestService) (*EditContestJSONRequestBody, string) {
				path := fmt.Sprintf("/api/v1/contests/%s", invalidID)
				return &EditContestJSONRequestBody{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long description",
			setup: func(_ *mock_service.MockContestService) (*EditContestJSONRequestBody, string) {
				contestID := random.UUID()
				description := strings.Repeat("a", 257)
				reqBody := &EditContestJSONRequestBody{
					Description: &description,
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid link",
			setup: func(_ *mock_service.MockContestService) (*EditContestJSONRequestBody, string) {
				contestID := random.UUID()
				link := random.AlphaNumeric()
				reqBody := &EditContestJSONRequestBody{
					Link: &link,
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid duration",
			setup: func(_ *mock_service.MockContestService) (*EditContestJSONRequestBody, string) {
				contestID := random.UUID()
				since, until := random.SinceAndUntil()
				since, until = until, since
				reqBody := &EditContestJSONRequestBody{
					Duration: &Duration{
						Since: since,
						Until: &until,
					},
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long name",
			setup: func(_ *mock_service.MockContestService) (*EditContestJSONRequestBody, string) {
				contestID := random.UUID()
				name := strings.Repeat("a", 33)
				reqBody := &EditContestJSONRequestBody{
					Name: &name,
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
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
				s.EXPECT().DeleteContest(anyCtx{}, contestID).Return(nil)
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
		setup      func(s *mock_service.MockContestService) (hres []*ContestTeam, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (hres []*ContestTeam, path string) {
				contestID := random.UUID()
				repoContestTeams := []*domain.ContestTeam{
					{
						ID:        random.UUID(),
						ContestID: contestID,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
						Members:   make([]*domain.User, 0),
					},
					{
						ID:        random.UUID(),
						ContestID: contestID,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
						Members:   make([]*domain.User, 0),
					},
				}
				hres = []*ContestTeam{
					{
						Id:      repoContestTeams[0].ID,
						Members: make([]User, 0),
						Name:    repoContestTeams[0].Name,
						Result:  repoContestTeams[0].Result,
					},
					{
						Id:      repoContestTeams[1].ID,
						Members: make([]User, 0),
						Name:    repoContestTeams[1].Name,
						Result:  repoContestTeams[1].Result,
					},
				}
				s.EXPECT().GetContestTeams(anyCtx{}, contestID).Return(repoContestTeams, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(s *mock_service.MockContestService) (hres []*ContestTeam, path string) {
				return []*ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", invalidID)
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

			hres := make([]*ContestTeam, 0, len(expectedHres))
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
		setup      func(s *mock_service.MockContestService) (hres ContestTeamDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (ContestTeamDetail, string) {
				teamID := random.UUID()
				contestID := random.UUID()
				repoContestTeamDetail := domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ID:        teamID,
						ContestID: contestID,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
						Members: []*domain.User{
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						},
					},
					Link:        random.AlphaNumeric(),
					Description: random.AlphaNumeric(),
				}
				members := make([]User, 0, len(repoContestTeamDetail.Members))
				for _, member := range repoContestTeamDetail.Members {
					members = append(members, User{
						Id:       member.ID,
						Name:     member.Name,
						RealName: member.RealName(),
					})
				}

				hres := ContestTeamDetail{
					Description: repoContestTeamDetail.Description,
					Id:          repoContestTeamDetail.ID,
					Link:        repoContestTeamDetail.Link,
					Members:     members,
					Name:        repoContestTeamDetail.Name,
					Result:      repoContestTeamDetail.Result,
				}

				s.EXPECT().GetContestTeam(anyCtx{}, contestID, teamID).Return(&repoContestTeamDetail, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(s *mock_service.MockContestService) (ContestTeamDetail, string) {
				return ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", invalidID, random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (ContestTeamDetail, string) {
				return ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "NotFound: Contest not found",
			setup: func(s *mock_service.MockContestService) (ContestTeamDetail, string) {
				teamID := random.UUID()
				contestID := random.UUID()
				s.EXPECT().GetContestTeam(anyCtx{}, contestID, teamID).Return(nil, repository.ErrNotFound)
				return ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
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

			var hres ContestTeamDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestContestHandler_AddContestTeam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (reqBody *AddContestTeamJSONRequestBody, expectedResBody ContestTeam, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
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
						Members:   make([]*domain.User, 0),
					},
					Link:        args.Link.String,
					Description: args.Description,
				}
				expectedResBody := ContestTeam{
					Id:      teamID,
					Members: make([]User, 0),
					Name:    want.Name,
					Result:  want.Result,
				}
				s.EXPECT().CreateContestTeam(anyCtx{}, contestID, &args).Return(&want, nil)
				return reqBody, expectedResBody, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(_ *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: missing required arg",
			setup: func(_ *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					// Name:        random.AlphaNumeric(), // missing
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long description",
			setup: func(_ *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: strings.Repeat("a", 257),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid link",
			setup: func(_ *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.AlphaNumeric()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long name",
			setup: func(_ *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        strings.Repeat("a", 33),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long result",
			setup: func(_ *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, strings.Repeat("a", 33)),
				}
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.StringFrom(reqBody.Result),
					Link:        optional.StringFrom(reqBody.Link),
					Description: reqBody.Description,
				}
				s.EXPECT().CreateContestTeam(anyCtx{}, contestID, &args).Return(nil, repository.ErrNotFound)
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "conflict contest",
			setup: func(s *mock_service.MockContestService) (*AddContestTeamJSONRequestBody, ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &AddContestTeamJSONRequestBody{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.StringFrom(reqBody.Result),
					Link:        optional.StringFrom(reqBody.Link),
					Description: reqBody.Description,
				}
				s.EXPECT().CreateContestTeam(anyCtx{}, contestID, &args).Return(nil, repository.ErrAlreadyExists)
				return reqBody, ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
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

			var resBody ContestTeam
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
		setup      func(s *mock_service.MockContestService) (reqBody *EditContestTeamJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumeric()),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumeric()),
					Description: ptr(t, random.AlphaNumeric()),
				}
				args := repository.UpdateContestTeamArgs{
					Name:        optional.StringFrom(reqBody.Name),
					Link:        optional.StringFrom(reqBody.Link),
					Result:      optional.StringFrom(reqBody.Result),
					Description: optional.StringFrom(reqBody.Description),
				}
				s.EXPECT().UpdateContestTeam(anyCtx{}, teamID, &args).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(_ *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				reqBody := &EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumeric()),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumeric()),
					Description: ptr(t, random.AlphaNumeric()),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", invalidID, random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(_ *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				reqBody := &EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumeric()),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumeric()),
					Description: ptr(t, random.AlphaNumeric()),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: not nil but empty",
			setup: func(_ *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				emptyStr := ""
				reqBody := &EditContestTeamJSONRequestBody{
					Description: &emptyStr,
					Name:        &emptyStr,
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: too long string",
			setup: func(_ *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				reqBody := &EditContestTeamJSONRequestBody{
					Description: ptr(t, strings.Repeat("a", 257)),
					Name:        ptr(t, strings.Repeat("a", 33)),
					Result:      ptr(t, strings.Repeat("a", 33)),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: invalid link",
			setup: func(_ *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				reqBody := &EditContestTeamJSONRequestBody{
					Link: ptr(t, random.AlphaNumeric()),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) (*EditContestTeamJSONRequestBody, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &EditContestTeamJSONRequestBody{
					Name:        ptr(t, random.AlphaNumeric()),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumeric()),
					Description: ptr(t, random.AlphaNumeric()),
				}
				args := repository.UpdateContestTeamArgs{
					Name:        optional.StringFrom(reqBody.Name),
					Link:        optional.StringFrom(reqBody.Link),
					Result:      optional.StringFrom(reqBody.Result),
					Description: optional.StringFrom(reqBody.Description),
				}
				s.EXPECT().UpdateContestTeam(anyCtx{}, teamID, &args).Return(repository.ErrNotFound)
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

func TestContestHandler_GetContestTeamMembers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockContestService) (hres []*User, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) ([]*User, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				users := []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				}
				hres := make([]*User, len(users))
				for i, user := range users {
					hres[i] = &User{
						Id:       user.ID,
						Name:     user.Name,
						RealName: user.RealName(),
					}
				}

				s.EXPECT().GetContestTeamMembers(anyCtx{}, contestID, teamID).Return(users, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(_ *mock_service.MockContestService) ([]*User, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(_ *mock_service.MockContestService) ([]*User, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) ([]*User, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				s.EXPECT().GetContestTeamMembers(anyCtx{}, contestID, teamID).Return(nil, repository.ErrNotFound)
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

			var hres []*User
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestContestHandler_AddContestTeamMembers(t *testing.T) {
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
				s.EXPECT().AddContestTeamMembers(anyCtx{}, teamID, reqBody.Members).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: members is empty",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				return &Req{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: memberID is invalid",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				return &Req{
					Members: []uuid.UUID{
						random.UUID(),
						uuid.Nil,
					},
				}, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
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
				s.EXPECT().AddContestTeamMembers(anyCtx{}, teamID, reqBody.Members).Return(repository.ErrNotFound)
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

func TestContestHandler_EditContestTeamMembers(t *testing.T) {
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
				s.EXPECT().EditContestTeamMembers(anyCtx{}, teamID, reqBody.Members).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: members is empty",
			setup: func(_ *mock_service.MockContestService) (*Req, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				return &Req{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
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
				s.EXPECT().EditContestTeamMembers(anyCtx{}, teamID, reqBody.Members).Return(repository.ErrNotFound)
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

			statusCode, _ := doRequest(t, api, http.MethodPut, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
