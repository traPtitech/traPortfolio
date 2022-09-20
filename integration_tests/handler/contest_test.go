package handler

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

// GetContests GET /contests
func TestGetContests(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockContests,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_get_contests")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContests), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetContest GET /contests/:contestID
func TestGetContest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.ContestID1(),
			mockdata.HMockContests[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("bad request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError("not found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_get_contest")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContest, tt.contestID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// CreateContest POST /contests
func TestCreateContest(t *testing.T) {
	var (
		name          = random.AlphaNumeric()
		link          = random.RandURLString()
		description   = random.AlphaNumeric()
		since, until  = random.SinceAndUntil()
		tooLongString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		invalidURL    = "invalid url"
		//tooLongStringは260文字
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqbody    handler.CreateContestJSONBody
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			handler.CreateContestJSONBody{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: name,
			},
			handler.ContestDetail{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Id:    uuid.Nil,
				Link:  link,
				Name:  name,
				Teams: []handler.ContestTeam{},
			},
		},
		"400 invalid description": {
			http.StatusBadRequest,
			handler.CreateContestJSONBody{
				Description: tooLongString,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: name,
			},
			testutils.HTTPError(repository.ErrValidate.Error()),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			handler.CreateContestJSONBody{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &invalidURL,
				Name: name,
			},
			testutils.HTTPError(repository.ErrValidate.Error()),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			handler.CreateContestJSONBody{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: tooLongString,
			},
			testutils.HTTPError(repository.ErrValidate.Error()),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_create_contests")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.Contest.CreateContest), &tt.reqbody)
			switch want := tt.want.(type) {
			case handler.ContestDetail:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res, testutils.OptSyncID, testutils.OptRetrieveID(&want.Id))
			case error:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

func TestEditContest(t *testing.T) {
}

func TestDeleteContest(t *testing.T) {
}

func TestGetContestTeam(t *testing.T) {
}

func TestAddContestTeam(t *testing.T) {
}

func TestEditContestTeam(t *testing.T) {
}

func TestGetContestTeamMember(t *testing.T) {
}

func TestAddContestTeamMember(t *testing.T) {
}

func TestEditContestTeamMember(t *testing.T) {
}

/*
// GetXXX GET /XXX
func TestGetXXX(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       interface{}
	}{
		// TODO: Add cases
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("xxx_handler_get_xxx")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
                        t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.XXX.GetXXX, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
*/
