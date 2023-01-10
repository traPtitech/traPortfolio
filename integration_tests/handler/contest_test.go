package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
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
			mockdata.CloneHandlerMockContestDetails()[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError("Not Found: not found"),
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
		tooLongString = strings.Repeat("a", 260)
		invalidURL    = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqbody    handler.CreateContestJSONRequestBody
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			handler.CreateContestJSONRequestBody{
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
			handler.CreateContestJSONRequestBody{
				Description: tooLongString,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: name,
			},
			testutils.HTTPError("Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			handler.CreateContestJSONRequestBody{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &invalidURL,
				Name: name,
			},
			testutils.HTTPError("Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			handler.CreateContestJSONRequestBody{
				Description: description,
				Duration: handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: tooLongString,
			},
			testutils.HTTPError("Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 since time is after until time": {
			http.StatusBadRequest,
			handler.CreateContestJSONRequestBody{
				Description: description,
				Duration: handler.Duration{
					Since: until,
					Until: &since,
				},
				Link: &link,
				Name: name,
			},
			testutils.HTTPError("Bad Request: validate error: duration: must be a valid date."),
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
	var (
		contest       = mockdata.CloneMockContests()[0]
		description   = contest.Description
		since         = contest.Since
		until         = contest.Until
		link          = contest.Link
		name          = contest.Name
		tooLongString = strings.Repeat("a", 260)
		invalidURL    = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		reqBody    handler.EditContestJSONRequestBody
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			handler.EditContestJSONRequestBody{
				Description: &description,
				Duration: &handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: &name,
			},
			nil,
		},
		"204 without change": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			handler.EditContestJSONRequestBody{},
			nil,
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.EditContestJSONRequestBody{},
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid description": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.EditContestJSONRequestBody{
				Description: &tooLongString,
			},
			testutils.HTTPError("Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.EditContestJSONRequestBody{
				Link: &invalidURL,
			},
			testutils.HTTPError("Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.EditContestJSONRequestBody{
				Name: &tooLongString,
			},
			testutils.HTTPError("Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 since time is after until time": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.EditContestJSONRequestBody{
				Duration: &handler.Duration{
					Since: until,
					Until: &since,
				},
			},
			testutils.HTTPError("Bad Request: validate error: duration: must be a valid date."),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			handler.EditContestJSONRequestBody{
				Description: &description,
				Duration: &handler.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: &name,
			},
			testutils.HTTPError("Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_edit_contest")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var contest handler.ContestDetail
				res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContest, tt.contestID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &contest)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない

				// Update & Assert
				res = testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.Contest.EditContest, tt.contestID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)

				// Get updated response & Assert
				res = testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContest, tt.contestID), nil)
				testutils.AssertResponse(t, http.StatusOK, contest, res)
			} else {
				res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.Contest.EditContest, tt.contestID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

func TestDeleteContest(t *testing.T) {
	// https://github.com/traPtitech/traPortfolio/issues/460
}

func TestGetContestTeam(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.ContestID1(),
			mockdata.CloneHandlerMockContestTeamsByID()[mockdata.ContestID1()],
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_get_contest_team")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeams, tt.contestID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

func TestAddContestTeam(t *testing.T) {
	var (
		description   = random.AlphaNumeric()
		link          = random.RandURLString()
		name          = random.AlphaNumeric()
		result        = random.AlphaNumeric()
		tooLongString = strings.Repeat("a", 260)
		invalidURL    = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		reqbody    handler.AddContestTeamJSONRequestBody
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			mockdata.ContestID1(),
			handler.AddContestTeamJSONRequestBody{
				Description: description,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			handler.ContestTeam{
				Id:     uuid.Nil, //テスト時にOptSyncIDで同期するため適当
				Name:   name,
				Result: result,
			},
		},
		"400 invalid description": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.AddContestTeamJSONRequestBody{
				Description: tooLongString,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			testutils.HTTPError("Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.AddContestTeamJSONRequestBody{
				Description: description,
				Link:        &invalidURL,
				Name:        name,
				Result:      &result,
			},
			testutils.HTTPError("Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			handler.AddContestTeamJSONRequestBody{
				Description: description,
				Link:        &link,
				Name:        tooLongString,
				Result:      &result,
			},
			testutils.HTTPError("Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		/*"404": {
			http.StatusNotFound,
			random.UUID(),
			handler.AddContestTeamJSONRequestBody{
				Description: description,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			testutils.HTTPError("Not Found: not found"),
		},*/
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_add_contest_team")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.Contest.AddContestTeam, tt.contestID), &tt.reqbody)
			switch tt.want.(type) {
			case handler.ContestDetail:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res, testutils.OptSyncID)
			case error:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

func TestEditContestTeam(t *testing.T) {
}

// GetContestTeamMembers GET /contests/:contestID/teams/:teamID/members
func TestGetContestTeamMembers(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		teamID     uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			[]handler.User{
				mockdata.CloneHandlerMockUsers()[0],
			},
		},
		"200 with no members": {
			http.StatusOK,
			mockdata.ContestID1(),
			mockdata.ContestTeamID2(),
			[]handler.User{},
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.ContestTeamID1(),
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			testutils.HTTPError("Bad Request: nil id"),
		},
		"404 contestID not exist": {
			http.StatusNotFound,
			random.UUID(),
			mockdata.ContestTeamID1(),
			testutils.HTTPError("Not Found: not found"),
		},
		"404 teamID not exist": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			testutils.HTTPError("Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_get_contest_team_members")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeamMembers, tt.contestID, tt.teamID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// AddContestTeamMembers POST /contests/:contestID/teams/:teamID/members
func TestAddContestTeamMembers(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		teamID     uuid.UUID
		reqbody    handler.AddContestTeamMembersJSONRequestBody
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			handler.AddContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			nil,
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.ContestTeamID1(),
			handler.AddContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			handler.AddContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid memberID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			handler.AddContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					uuid.Nil,
				},
			},
			testutils.HTTPError("Bad Request: validate error: members: (0: must be a valid UUID v4.)."),
		},
		"400 invalid member": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			handler.AddContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					random.UUID(),
				},
			},
			testutils.HTTPError("Bad Request: argument error"),
		},
		"404 team not found": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			handler.AddContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			testutils.HTTPError("Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_add_contest_team_member")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.Contest.AddContestTeamMembers, tt.contestID, tt.teamID), &tt.reqbody)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// EditContestTeamMembers PUT /contests/:contestID/teams/:teamID/members
func TestEditContestTeamMembers(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		teamID     uuid.UUID
		reqbody    handler.EditContestTeamMembersJSONRequestBody
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			handler.EditContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			nil,
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.ContestTeamID1(),
			handler.EditContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			handler.EditContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid memberID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			handler.EditContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					uuid.Nil,
				},
			},
			testutils.HTTPError("Bad Request: validate error: members: (0: must be a valid UUID v4.)."),
		},
		"400 invalid member": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			handler.EditContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					random.UUID(),
				},
			},
			testutils.HTTPError("Bad Request: argument error"),
		},
		"404 team not found": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			handler.EditContestTeamMembersJSONRequestBody{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			testutils.HTTPError("Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("contest_handler_edit_contest_team_member")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Update & Assert
				res := testutils.DoRequest(t, e, http.MethodPut, e.URL(api.Contest.EditContestTeamMembers, tt.contestID, tt.teamID), &tt.reqbody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)

				// Assert
				res = testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeamMembers, tt.contestID, tt.teamID), nil)
				var response []model.User
				var userIDs []uuid.UUID
				err := json.Unmarshal(res.Body.Bytes(), &response)
				if err != nil {
					assert.Error(t, err)
				}
				for _, memberID := range response {
					userIDs = append(userIDs, memberID.ID)
				}
				assert.Equal(t, tt.reqbody.Members, userIDs)
			} else {
				res := testutils.DoRequest(t, e, http.MethodPut, e.URL(api.Contest.EditContestTeamMembers, tt.contestID, tt.teamID), &tt.reqbody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
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
