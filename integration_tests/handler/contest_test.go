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
func GetContests(t *testing.T) {
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
func GetContest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockContest.Id,
			mockdata.HMockContest,
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
func CreateContest(t *testing.T) {
	var (
		name         = random.AlphaNumeric()
		link         = random.RandURLString()
		description  = random.AlphaNumeric()
		since, until = random.SinceAndUntil()
		invalidUrl   = "invalid url"
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
		"400 invalid Link": {
			http.StatusBadRequest,
			handler.CreateContestJSONBody{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &invalidUrl,
				Name: name,
			},
			testutils.HTTPError(repository.ErrValidate.Error()),
		},
		// TODO: validationもテストする
		// https://github.com/traPtitech/traPortfolio/pull/391#discussion_r952794355
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

func EditContest(t *testing.T) {
}

func DeleteContest(t *testing.T) {
}

func GetContestTeam(t *testing.T) {
}

func AddContestTeam(t *testing.T) {
}

func EditContestTeam(t *testing.T) {
}

func GetContestTeamMember(t *testing.T) {
}

func AddContestTeamMember(t *testing.T) {
}

func EditContestTeamMember(t *testing.T) {
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
