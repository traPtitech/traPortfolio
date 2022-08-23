package handler

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

// GetContests GET /contests
func TestContestHandler_GetContests(t *testing.T) {
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
func TestContestHandler_GetContest(t *testing.T) {
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
