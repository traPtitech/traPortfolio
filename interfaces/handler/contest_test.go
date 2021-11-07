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
	"github.com/traPtitech/traPortfolio/util/random"
)

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestContestHandler_GetContests(t *testing.T) {
	until := mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00")

	tests := []struct {
		name         string
		setup        func(th *handler.TestHandlers, want []*domain.Contest) (path string)
		statusCode   int
		repoContests []*domain.Contest
		hresContests []*handler.Contest
	}{
		{
			name: "success",
			setup: func(th *handler.TestHandlers, want []*domain.Contest) string {
				th.Service.MockContestService.EXPECT().GetContests(gomock.Any()).Return(want, nil)
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
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			for i, v := range tt.hresContests {
				tt.repoContests[i].ID = v.Id
			}
			path := tt.setup(&handlers, tt.repoContests)

			var resBody []*handler.Contest
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, path, nil, &resBody)

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

func makeContest() (*domain.ContestDetail, *handler.ContestDetail) {
	d := domain.ContestDetail{
		Contest: domain.Contest{
			ID:        getContestID[0],
			Name:      random.AlphaNumeric(rand.Intn(30) + 1),
			TimeStart: random.Time(),
			TimeEnd:   random.Time(),
		},
		Link:        random.RandURLString(),
		Description: random.AlphaNumeric(rand.Intn(30) + 1),
		Teams: []*domain.ContestTeam{
			{
				ID:        getContestID[1],
				ContestID: getContestID[0],
				Name:      random.AlphaNumeric(rand.Intn(30) + 1),
				Result:    random.AlphaNumeric(rand.Intn(30) + 1),
			},
			{
				ID:        getContestID[2],
				ContestID: getContestID[0],
				Name:      random.AlphaNumeric(rand.Intn(30) + 1),
				Result:    random.AlphaNumeric(rand.Intn(30) + 1),
			},
		},
	}

	teams := make([]handler.ContestTeam, len(d.Teams))
	for i, v := range d.Teams {
		teams[i] = handler.ContestTeam{
			Id:     v.ID,
			Name:   v.Name,
			Result: &v.Result,
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
		Link:        &d.Link,
		Description: d.Description,
		Teams:       teams,
	}

	return &d, &hres
}

func TestContestHandler_GetContest(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (repoContest *domain.ContestDetail, hresContest *handler.ContestDetail, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(th *handler.TestHandlers) (*domain.ContestDetail, *handler.ContestDetail, string) {
				want, hres := makeContest()
				th.Service.MockContestService.EXPECT().GetContest(gomock.Any(), want.ID).Return(want, nil)
				path := fmt.Sprintf("/api/v1/contests/%s", want.ID.String())

				return want, hres, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			setup: func(th *handler.TestHandlers) (*domain.ContestDetail, *handler.ContestDetail, string) {
				path := "/api/v1/contests/invalid"
				return &domain.ContestDetail{}, &handler.ContestDetail{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			setup: func(th *handler.TestHandlers) (*domain.ContestDetail, *handler.ContestDetail, string) {
				uid := random.UUID()
				th.Service.MockContestService.EXPECT().GetContest(gomock.Any(), uid).Return(nil, repository.ErrNotFound)

				return &domain.ContestDetail{}, &handler.ContestDetail{}, fmt.Sprintf("/api/v1/contests/%s", uid)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			_, hresContest, path := tt.setup(&handlers)

			var resBody handler.ContestDetail
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresContest, &resBody)
		})
	}
}

// func TestContestHandler_PostContest(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setup      func(th *handler.TestHandlers) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string)
// 		statusCode int
// 	}{
// 		{
// 			name: "Success",
// 			setup: func(th *handler.TestHandlers) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string) {
// 				reqBody = &handler.PostContestJSONRequestBody{
// 					Description: random.AlphaNumeric(rand.Intn(30) + 1),
// 					Duration: handler.Duration{
// 						Since: random.Time(),
// 						Until: random.Time(),
// 					},
// 					Name: random.AlphaNumeric(rand.Intn(30) + 1),
// 					Link: random.RandURLString(),
// 				}
// 				args := repository.CreateContestArgs{
// 					Name:        reqBody.Name,
// 					Description: reqBody.Description,
// 					Link:        reqBody.Link,
// 					Since:       reqBody.Duration.Since,
// 					Until:       reqBody.Duration.Until,
// 				}
// 				want := domain.Contest{
// 					ID:        random.UUID(),
// 					Name:      args.Name,
// 					TimeStart: args.Since,
// 					TimeEnd:   args.Until,
// 				}
// 				expectedResBody = &handler.Contest{
// 					Id:   want.ID,
// 					Name: want.Name,
// 					Duration: handler.Duration{
// 						Since: want.TimeStart,
// 						Until: &want.TimeEnd,
// 					},
// 				}
// 				th.Service.MockContestService.EXPECT().CreateContest(gomock.Any(), &args).Return(&want, nil)
// 				path = "/api/v1/contests"
// 				return reqBody, expectedResBody, &handler.Contest{}, path
// 			},
// 			statusCode: http.StatusCreated,
// 		},
// 		{
// 			name: "Bad Request: invalid url",
// 			setup: func(th *handler.TestHandlers) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string) {
// 				reqBody = &handler.PostContestJSONRequestBody{
// 					Description: random.AlphaNumeric(rand.Intn(30) + 1),
// 					Duration: handler.Duration{
// 						Since: random.Time(),
// 						Until: random.Time(),
// 					},
// 					Name: random.AlphaNumeric(rand.Intn(30) + 1),
// 					Link: random.AlphaNumeric(rand.Intn(30) + 1),
// 				}
// 				path = "/api/v1/contests"
// 				return reqBody, nil, nil, path
// 			},
// 			statusCode: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Conflict",
// 			setup: func(th *handler.TestHandlers) (reqBody *handler.PostContestJSONRequestBody, expectedResBody *handler.Contest, resBody *handler.Contest, path string) {
// 				reqBody = &handler.PostContestJSONRequestBody{
// 					Description: random.AlphaNumeric(rand.Intn(30) + 1),
// 					Duration: handler.Duration{
// 						Since: random.Time(),
// 						Until: random.Time(),
// 					},
// 					Name: random.AlphaNumeric(rand.Intn(30) + 1),
// 					Link: random.RandURLString(),
// 				}
// 				args := repository.CreateContestArgs{
// 					Name:        reqBody.Name,
// 					Description: reqBody.Description,
// 					Link:        reqBody.Link,
// 					Since:       reqBody.Duration.Since,
// 					Until:       reqBody.Duration.Until,
// 				}
// 				th.Service.MockContestService.EXPECT().CreateContest(gomock.Any(), &args).Return(nil, repository.ErrAlreadyExists)
// 				return reqBody, nil, nil, "/api/v1/contests"
// 			},
// 			statusCode: http.StatusConflict,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			handlers := SetupTestHandlers(t, ctrl)

// 			reqBody, res, resBody, path := tt.setup(&handlers)

// 			statusCode, _ := doRequest(t, handlers.API, http.MethodPost, path, reqBody, resBody)

// 			// Assertion
// 			assert.Equal(t, tt.statusCode, statusCode)
// 			assert.Equal(t, res, resBody)
// 		})
// 	}
// }

// func TestContestHandler_PatchContest(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.PatchContest(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_DeleteContest(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.DeleteContest(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_GetContestTeams(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.GetContestTeams(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_GetContestTeam(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.GetContestTeam(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_PostContestTeam(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.PostContestTeam(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_PatchContestTeam(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.PatchContestTeam(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_GetContestTeamMember(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.GetContestTeamMember(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_PostContestTeamMember(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.PostContestTeamMember(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_DeleteContestTeamMember(t *testing.T) {
// 	type fields struct {
// 		srv service.ContestService
// 	}
// 	type args struct {
// 		_c echo.Context
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		setup     func(f fields, args args)
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mock
// 			ctrl := gomock.NewController(t)
// 			tt.fields = fields{
// 				srv: mock_service.NewMockContestService(ctrl),
// 			}
// 			tt.setup(tt.fields, tt.args)
// 			h := NewContestHandler(tt.fields.srv)
// 			// Assertion
// 			tt.assertion(t, h.DeleteContestTeamMember(tt.args._c))
// 		})
// 	}
// }
