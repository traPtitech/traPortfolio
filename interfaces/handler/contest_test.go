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
	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
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

func mustParseTime(t *testing.T, layout, value string) time.Time {
	t.Helper()

	tm, err := time.Parse(layout, value)
	assert.NoError(t, err)
	return tm
}

func TestContestHandler_GetContests(t *testing.T) {
	until := mustParseTime(t, time.RFC3339, "2006-01-02T15:04:05+09:00")

	t.Parallel()
	tests := []struct {
		name         string
		setup        func(s *mock_service.MockContestService, want []*domain.Contest) (path string)
		statusCode   int
		repoContests []*domain.Contest
		hresContests []*schema.Contest
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
					TimeStart: mustParseTime(t, time.RFC3339, "2006-01-02T15:04:05+09:00"),
					TimeEnd:   mustParseTime(t, time.RFC3339, "2006-01-02T15:04:05+09:00"),
				},
			},
			hresContests: []*schema.Contest{
				{
					Name: "test1",
					Duration: schema.Duration{
						Since: mustParseTime(t, time.RFC3339, "2006-01-02T15:04:05+09:00"),
						Until: &until,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			for i, v := range tt.hresContests {
				tt.repoContests[i].ID = v.Id
			}
			path := tt.setup(s, tt.repoContests)

			var resBody []*schema.Contest
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

func makeContest(t *testing.T) (*domain.ContestDetail, *schema.ContestDetail) {
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
				ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
					ID:        getContestID[1],
					ContestID: getContestID[0],
					Name:      random.AlphaNumeric(),
					Result:    random.AlphaNumeric(),
				},
				Members: []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				},
			},
			{
				ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
					ID:        getContestID[2],
					ContestID: getContestID[0],
					Name:      random.AlphaNumeric(),
					Result:    random.AlphaNumeric(),
				},
				Members: []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				},
			},
		},
	}

	teams := make([]schema.ContestTeam, len(d.ContestTeams))
	for i, v := range d.ContestTeams {
		member := make([]schema.User, len(v.Members))
		for j, w := range v.Members {
			member[j] = schema.User{
				Id:       w.ID,
				Name:     w.Name,
				RealName: w.RealName(),
			}
		}
		teams[i] = schema.ContestTeam{
			Id:      v.ID,
			Members: member,
			Name:    v.Name,
			Result:  v.Result,
		}
	}

	hres := schema.ContestDetail{
		Description: d.Description,
		Duration: schema.Duration{
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
		setup      func(s *mock_service.MockContestService) (repoContest *domain.ContestDetail, hresContest *schema.ContestDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *schema.ContestDetail, string) {
				want, hres := makeContest(t)
				s.EXPECT().GetContest(anyCtx{}, want.ID).Return(want, nil)
				path := fmt.Sprintf("/api/v1/contests/%s", want.ID.String())

				return want, hres, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			setup: func(_ *mock_service.MockContestService) (*domain.ContestDetail, *schema.ContestDetail, string) {
				path := fmt.Sprintf("/api/v1/contests/%s", invalidID)
				return &domain.ContestDetail{}, &schema.ContestDetail{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockContestService) (*domain.ContestDetail, *schema.ContestDetail, string) {
				uid := random.UUID()
				s.EXPECT().GetContest(anyCtx{}, uid).Return(nil, repository.ErrNotFound)

				return &domain.ContestDetail{}, &schema.ContestDetail{}, fmt.Sprintf("/api/v1/contests/%s", uid)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			_, hresContest, path := tt.setup(s)

			var resBody schema.ContestDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresContest, &resBody)
		})
	}
}

func makeCreateContestRequest(t *testing.T, description string, since time.Time, until time.Time, name string, link string) *schema.CreateContestRequest {
	t.Helper()
	return &schema.CreateContestRequest{
		Description: description,
		Duration: schema.Duration{
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
		setup      func(s *mock_service.MockContestService) (reqBody *schema.CreateContestRequest, expectedResBody *schema.Contest, resBody *schema.Contest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (reqBody *schema.CreateContestRequest, expectedResBody *schema.Contest, resBody *schema.Contest, path string) {
				since, until := random.SinceAndUntil()
				reqBody = makeCreateContestRequest(
					t,
					random.AlphaNumeric(),
					since,
					until,
					random.AlphaNumeric(),
					random.RandURLString(),
				)
				args := repository.CreateContestArgs{
					Name:        reqBody.Name,
					Description: reqBody.Description,
					Link:        optional.FromPtr(reqBody.Link),
					Since:       reqBody.Duration.Since,
					Until:       optional.FromPtr(reqBody.Duration.Until),
				}
				want := domain.ContestDetail{
					Contest: domain.Contest{
						ID:        random.UUID(),
						Name:      args.Name,
						TimeStart: args.Since,
						TimeEnd:   args.Until.ValueOrZero(),
					},
					Link:         args.Link.ValueOrZero(),
					Description:  args.Description,
					ContestTeams: []*domain.ContestTeam{},
				}
				expectedResBody = &schema.Contest{
					Id:   want.ID,
					Name: want.Name,
					Duration: schema.Duration{
						Since: want.TimeStart,
						Until: &want.TimeEnd,
					},
				}
				s.EXPECT().CreateContest(anyCtx{}, &args).Return(&want, nil)
				path = "/api/v1/contests"
				return reqBody, expectedResBody, &schema.Contest{}, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Bad Request: invalid url",
			setup: func(_ *mock_service.MockContestService) (reqBody *schema.CreateContestRequest, expectedResBody *schema.Contest, resBody *schema.Contest, path string) {
				since, until := random.SinceAndUntil()
				reqBody = makeCreateContestRequest(
					t,
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
			setup: func(s *mock_service.MockContestService) (reqBody *schema.CreateContestRequest, expectedResBody *schema.Contest, resBody *schema.Contest, path string) {
				since, until := random.SinceAndUntil()
				reqBody = makeCreateContestRequest(
					t,
					random.AlphaNumeric(),
					since,
					until,
					random.AlphaNumeric(),
					random.RandURLString(),
				)
				args := repository.CreateContestArgs{
					Name:        reqBody.Name,
					Description: reqBody.Description,
					Link:        optional.FromPtr(reqBody.Link),
					Since:       reqBody.Duration.Since,
					Until:       optional.FromPtr(reqBody.Duration.Until),
				}
				s.EXPECT().CreateContest(anyCtx{}, &args).Return(nil, repository.ErrAlreadyExists)
				return reqBody, nil, nil, "/api/v1/contests"
			},
			statusCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
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
		setup      func(s *mock_service.MockContestService) (reqBody *schema.EditContestRequest, path string)
		statusCode int
	}{
		{
			name: "Success 1",
			setup: func(s *mock_service.MockContestService) (*schema.EditContestRequest, string) {
				contestID := random.UUID()
				name := random.AlphaNumeric()
				link := random.RandURLString()
				description := random.AlphaNumeric()
				since, until := random.SinceAndUntil()
				reqBody := &schema.EditContestRequest{
					Name:        &name,
					Link:        &link,
					Description: &description,
					Duration: &schema.Duration{
						Since: since,
						Until: &until,
					},
				}
				args := repository.UpdateContestArgs{
					Name:        optional.FromPtr(reqBody.Name),
					Description: optional.FromPtr(reqBody.Description),
					Link:        optional.FromPtr(reqBody.Link),
					Since:       optional.FromPtr(&reqBody.Duration.Since),
					Until:       optional.FromPtr(reqBody.Duration.Until),
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				s.EXPECT().UpdateContest(anyCtx{}, contestID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestRequest, string) {
				path := fmt.Sprintf("/api/v1/contests/%s", invalidID)
				return &schema.EditContestRequest{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long description",
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestRequest, string) {
				contestID := random.UUID()
				description := strings.Repeat("a", 257)
				reqBody := &schema.EditContestRequest{
					Description: &description,
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid link",
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestRequest, string) {
				contestID := random.UUID()
				link := random.AlphaNumeric()
				reqBody := &schema.EditContestRequest{
					Link: &link,
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid duration",
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestRequest, string) {
				contestID := random.UUID()
				since, until := random.SinceAndUntil()
				since, until = until, since
				reqBody := &schema.EditContestRequest{
					Duration: &schema.Duration{
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
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestRequest, string) {
				contestID := random.UUID()
				name := strings.Repeat("a", 33)
				reqBody := &schema.EditContestRequest{
					Name: &name,
				}
				path := fmt.Sprintf("/api/v1/contests/%s", contestID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
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
		setup      func(s *mock_service.MockContestService) (hres []*schema.ContestTeam, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (hres []*schema.ContestTeam, path string) {
				contestID := random.UUID()
				repoContestTeams := []*domain.ContestTeam{
					{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        random.UUID(),
							ContestID: contestID,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
						Members: []*domain.User{
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						},
					},
					{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        random.UUID(),
							ContestID: contestID,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
						Members: []*domain.User{
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						},
					},
				}
				hres = []*schema.ContestTeam{
					{
						Id: repoContestTeams[0].ID,
						Members: []schema.User{
							{Id: repoContestTeams[0].Members[0].ID, Name: repoContestTeams[0].Members[0].Name, RealName: repoContestTeams[0].Members[0].RealName()},
							{Id: repoContestTeams[0].Members[1].ID, Name: repoContestTeams[0].Members[1].Name, RealName: repoContestTeams[0].Members[1].RealName()},
						},
						Name:   repoContestTeams[0].Name,
						Result: repoContestTeams[0].Result,
					},
					{
						Id: repoContestTeams[1].ID,
						Members: []schema.User{
							{Id: repoContestTeams[1].Members[0].ID, Name: repoContestTeams[1].Members[0].Name, RealName: repoContestTeams[1].Members[0].RealName()},
							{Id: repoContestTeams[1].Members[1].ID, Name: repoContestTeams[1].Members[1].Name, RealName: repoContestTeams[1].Members[1].RealName()},
						},
						Name:   repoContestTeams[1].Name,
						Result: repoContestTeams[1].Result,
					},
				}
				s.EXPECT().GetContestTeams(anyCtx{}, contestID).Return(repoContestTeams, nil)
				return hres, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid ID",
			setup: func(s *mock_service.MockContestService) (hres []*schema.ContestTeam, path string) {
				return []*schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			expectedHres, path := tt.setup(s)

			hres := make([]*schema.ContestTeam, 0, len(expectedHres))
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
		setup      func(s *mock_service.MockContestService) (hres schema.ContestTeamDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (schema.ContestTeamDetail, string) {
				teamID := random.UUID()
				contestID := random.UUID()
				repoContestTeamDetail := domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        teamID,
							ContestID: contestID,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
						Members: []*domain.User{
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						},
					},
					Link:        random.AlphaNumeric(),
					Description: random.AlphaNumeric(),
				}
				members := make([]schema.User, 0, len(repoContestTeamDetail.Members))
				for _, member := range repoContestTeamDetail.Members {
					members = append(members, schema.User{
						Id:       member.ID,
						Name:     member.Name,
						RealName: member.RealName(),
					})
				}

				hres := schema.ContestTeamDetail{
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
			setup: func(s *mock_service.MockContestService) (schema.ContestTeamDetail, string) {
				return schema.ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", invalidID, random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(s *mock_service.MockContestService) (schema.ContestTeamDetail, string) {
				return schema.ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "NotFound: Contest not found",
			setup: func(s *mock_service.MockContestService) (schema.ContestTeamDetail, string) {
				teamID := random.UUID()
				contestID := random.UUID()
				s.EXPECT().GetContestTeam(anyCtx{}, contestID, teamID).Return(nil, repository.ErrNotFound)
				return schema.ContestTeamDetail{}, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			expectedHres, path := tt.setup(s)

			var hres schema.ContestTeamDetail
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
		setup      func(s *mock_service.MockContestService) (reqBody *schema.AddContestTeamRequest, expectedResBody schema.ContestTeam, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.FromPtr(reqBody.Result),
					Link:        optional.FromPtr(reqBody.Link),
					Description: reqBody.Description,
				}
				want := domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        teamID,
							ContestID: contestID,
							Name:      args.Name,
							Result:    args.Result.ValueOrZero(),
						},
						Members: make([]*domain.User, 0),
					},
					Link:        args.Link.ValueOrZero(),
					Description: args.Description,
				}
				expectedResBody := schema.ContestTeam{
					Id:      teamID,
					Members: make([]schema.User, 0),
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
			setup: func(_ *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: missing required arg",
			setup: func(_ *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					// Name:        random.AlphaNumeric(), // missing
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long description",
			setup: func(_ *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: strings.Repeat("a", 257),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid link",
			setup: func(_ *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.AlphaNumeric()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long name",
			setup: func(_ *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        strings.Repeat("a", 33),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: too long result",
			setup: func(_ *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, strings.Repeat("a", 33)),
				}
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.FromPtr(reqBody.Result),
					Link:        optional.FromPtr(reqBody.Link),
					Description: reqBody.Description,
				}
				s.EXPECT().CreateContestTeam(anyCtx{}, contestID, &args).Return(nil, repository.ErrNotFound)
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "conflict contest",
			setup: func(s *mock_service.MockContestService) (*schema.AddContestTeamRequest, schema.ContestTeam, string) {
				contestID := random.UUID()
				reqBody := &schema.AddContestTeamRequest{
					Name:        random.AlphaNumeric(),
					Link:        ptr(t, random.RandURLString()),
					Description: random.AlphaNumeric(),
					Result:      ptr(t, random.AlphaNumeric()),
				}
				args := repository.CreateContestTeamArgs{
					Name:        reqBody.Name,
					Result:      optional.FromPtr(reqBody.Result),
					Link:        optional.FromPtr(reqBody.Link),
					Description: reqBody.Description,
				}
				s.EXPECT().CreateContestTeam(anyCtx{}, contestID, &args).Return(nil, repository.ErrAlreadyExists)
				return reqBody, schema.ContestTeam{}, fmt.Sprintf("/api/v1/contests/%s/teams", contestID)
			},
			statusCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			reqBody, res, path := tt.setup(s)

			var resBody schema.ContestTeam
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
		setup      func(s *mock_service.MockContestService) (reqBody *schema.EditContestTeamRequest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &schema.EditContestTeamRequest{
					Name:        ptr(t, random.AlphaNumeric()),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumeric()),
					Description: ptr(t, random.AlphaNumeric()),
				}
				args := repository.UpdateContestTeamArgs{
					Name:        optional.FromPtr(reqBody.Name),
					Link:        optional.FromPtr(reqBody.Link),
					Result:      optional.FromPtr(reqBody.Result),
					Description: optional.FromPtr(reqBody.Description),
				}
				s.EXPECT().UpdateContestTeam(anyCtx{}, teamID, &args).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid contest ID",
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				reqBody := &schema.EditContestTeamRequest{
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
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				reqBody := &schema.EditContestTeamRequest{
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
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				emptyStr := ""
				reqBody := &schema.EditContestTeamRequest{
					Description: &emptyStr,
					Name:        &emptyStr,
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid request body: too long string",
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				reqBody := &schema.EditContestTeamRequest{
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
			setup: func(_ *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				reqBody := &schema.EditContestTeamRequest{
					Link: ptr(t, random.AlphaNumeric()),
				}
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", random.UUID(), random.UUID())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) (*schema.EditContestTeamRequest, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				reqBody := &schema.EditContestTeamRequest{
					Name:        ptr(t, random.AlphaNumeric()),
					Link:        ptr(t, random.RandURLString()),
					Result:      ptr(t, random.AlphaNumeric()),
					Description: ptr(t, random.AlphaNumeric()),
				}
				args := repository.UpdateContestTeamArgs{
					Name:        optional.FromPtr(reqBody.Name),
					Link:        optional.FromPtr(reqBody.Link),
					Result:      optional.FromPtr(reqBody.Result),
					Description: optional.FromPtr(reqBody.Description),
				}
				s.EXPECT().UpdateContestTeam(anyCtx{}, teamID, &args).Return(repository.ErrNotFound)
				return reqBody, fmt.Sprintf("/api/v1/contests/%s/teams/%s", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
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
		setup      func(s *mock_service.MockContestService) (hres []*schema.User, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockContestService) ([]*schema.User, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				users := []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				}
				hres := make([]*schema.User, len(users))
				for i, user := range users {
					hres[i] = &schema.User{
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
			setup: func(_ *mock_service.MockContestService) ([]*schema.User, string) {
				teamID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", invalidID, teamID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: Invalid team ID",
			setup: func(_ *mock_service.MockContestService) ([]*schema.User, string) {
				contestID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Contest not exist",
			setup: func(s *mock_service.MockContestService) ([]*schema.User, string) {
				contestID := random.UUID()
				teamID := random.UUID()
				s.EXPECT().GetContestTeamMembers(anyCtx{}, contestID, teamID).Return(nil, repository.ErrNotFound)
				return nil, fmt.Sprintf("/api/v1/contests/%s/teams/%s/members", contestID, teamID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupContestMock(t)

			expectedHres, path := tt.setup(s)

			var hres []*schema.User
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
