package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/traPtitech/traPortfolio/domain"
)

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestContestHandler_GetContests(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(th *TestHandlers, want []*domain.Contest) (handler echo.HandlerFunc, path string)
		statusCode   int
		dbContest    []*domain.Contest
		expectedBody []*Contest
		assertion    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			setup: func(th *TestHandlers, want []*domain.Contest) (handler echo.HandlerFunc, path string) {
				th.Repository.MockContestRepository.EXPECT().GetContests().Return(want, nil)
				return th.API.Contest.GetContests, "/api/v1/contests"
			},
			statusCode: 200,
			dbContest: []*domain.Contest{
				{
					ID:        uuid.Nil,
					Name:      "test1",
					TimeStart: mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"),
					TimeEnd:   mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"),
				},
			},
			expectedBody: []*Contest{
				{
					Name: "test1",
					Duration: Duration{
						Since: mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"),
						// Until: mustParseTime(time.RFC3339, "2006-01-02T15:04:05+09:00"), //TODO
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			for i, v := range tt.expectedBody {
				tt.dbContest[i].ID = v.Id
			}
			handler, path := tt.setup(&handlers, tt.dbContest)

			var resBody []*Contest
			statusCode, _, err := doRequest(t, handler, http.MethodGet, path, nil, &resBody)

			// Assertion
			tt.assertion(t, err)
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, tt.expectedBody, resBody)
		})
	}
}

// func TestContestHandler_GetContest(t *testing.T) {
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
// 			tt.assertion(t, h.GetContest(tt.args._c))
// 		})
// 	}
// }

// func TestContestHandler_PostContest(t *testing.T) {
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
// 			tt.assertion(t, h.PostContest(tt.args._c))
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
