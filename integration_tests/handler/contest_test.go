package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure/repository/model"
	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
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
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContests), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
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
			httpError(t, "Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContest, tt.contestID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// CreateContest POST /contests
func TestCreateContest(t *testing.T) {
	var (
		name                    = random.AlphaNumeric()
		link                    = random.RandURLString()
		description             = random.AlphaNumeric()
		since, until            = random.SinceAndUntil()
		tooLongString           = strings.Repeat("a", 260)
		justCountDescription    = strings.Repeat("亜", 256)
		justCountName           = strings.Repeat("亜", 32)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		invalidURL              = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqbody    schema.CreateContestRequest
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			schema.CreateContestRequest{
				Description: description,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: name,
			},
			schema.ContestDetail{
				Description: description,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Id:    uuid.Nil,
				Link:  link,
				Name:  name,
				Teams: []schema.ContestTeam{},
			},
		},
		"201 with Kanji": {
			http.StatusCreated,
			schema.CreateContestRequest{
				Description: justCountDescription,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: justCountName,
			},
			schema.ContestDetail{
				Description: justCountDescription,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Id:    uuid.Nil,
				Link:  link,
				Name:  justCountName,
				Teams: []schema.ContestTeam{},
			},
		},
		"400 invalid description": {
			http.StatusBadRequest,
			schema.CreateContestRequest{
				Description: tooLongString,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: name,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid description with Kanji": {
			http.StatusBadRequest,
			schema.CreateContestRequest{
				Description: tooLongDescriptionKanji,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: name,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			schema.CreateContestRequest{
				Description: description,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &invalidURL,
				Name: name,
			},
			httpError(t, "Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			schema.CreateContestRequest{
				Description: description,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: tooLongString,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Name with Kanji": {
			http.StatusBadRequest,
			schema.CreateContestRequest{
				Description: description,
				Duration: schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: tooLongName,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 since time is after until time": {
			http.StatusBadRequest,
			schema.CreateContestRequest{
				Description: description,
				Duration: schema.Duration{
					Since: until,
					Until: &since,
				},
				Link: &link,
				Name: name,
			},
			httpError(t, "Bad Request: validate error: duration: must be a valid date."),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodPost, e.URL(api.Contest.CreateContest), &tt.reqbody)
			switch want := tt.want.(type) {
			case schema.ContestDetail:
				assertResponse(t, tt.statusCode, tt.want, res, optSyncID, optRetrieveID(&want.Id))
			case error:
				assertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// EditContest PATCH /contests/:contestID
func TestEditContest(t *testing.T) {
	var (
		description             = random.AlphaNumeric()
		since, until            = random.SinceAndUntil()
		link                    = random.RandURLString()
		name                    = random.AlphaNumeric()
		tooLongString           = strings.Repeat("a", 260)
		justCountDescription    = strings.Repeat("亜", 256)
		justCountName           = strings.Repeat("亜", 32)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		invalidURL              = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		reqBody    schema.EditContestRequest
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Description: &description,
				Duration: &schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: &name,
			},
			nil,
		},
		"204 with kanji": {
			http.StatusNoContent,
			mockdata.ContestID2(),
			schema.EditContestRequest{
				Description: &justCountDescription,
				Duration: &schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: &justCountName,
			},
			nil,
		},
		"204 without change": {
			http.StatusNoContent,
			mockdata.ContestID3(),
			schema.EditContestRequest{},
			nil,
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			schema.EditContestRequest{},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid description": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Description: &tooLongString,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid description with kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Description: &tooLongDescriptionKanji,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Link: &invalidURL,
			},
			httpError(t, "Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Name: &tooLongString,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Name with kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Name: &tooLongName,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 since time is after until time": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.EditContestRequest{
				Duration: &schema.Duration{
					Since: until,
					Until: &since,
				},
			},
			httpError(t, "Bad Request: validate error: duration: must be a valid date."),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			schema.EditContestRequest{
				Description: &description,
				Duration: &schema.Duration{
					Since: since,
					Until: &until,
				},
				Link: &link,
				Name: &name,
			},
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var contest schema.ContestDetail
				res := doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContest, tt.contestID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &contest)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない

				// Update & Assert
				res = doRequest(t, e, http.MethodPatch, e.URL(api.Contest.EditContest, tt.contestID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)

				// Get updated response & Assert
				if tt.reqBody.Description != nil {
					contest.Description = *tt.reqBody.Description
				}
				if tt.reqBody.Duration != nil {
					contest.Duration = *tt.reqBody.Duration
				}
				if tt.reqBody.Link != nil {
					contest.Link = *tt.reqBody.Link
				}
				if tt.reqBody.Name != nil {
					contest.Name = *tt.reqBody.Name
				}
				res = doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContest, tt.contestID), nil)
				assertResponse(t, http.StatusOK, contest, res)
			} else {
				res := doRequest(t, e, http.MethodPatch, e.URL(api.Contest.EditContest, tt.contestID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// DeleteContest DELETE /contests/:contestID
func TestDeleteContest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			nil,
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodDelete, e.URL(api.Contest.DeleteContest, tt.contestID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetContestTeams GET /contests/:contestID/teams
func TestGetContestTeams(t *testing.T) {
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
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeams, tt.contestID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// AddContestTeam POST /contests/:contestID/teams
func TestAddContestTeam(t *testing.T) {
	var (
		description             = random.AlphaNumeric()
		link                    = random.RandURLString()
		name                    = random.AlphaNumeric()
		result                  = random.AlphaNumeric()
		tooLongString           = strings.Repeat("a", 260)
		justCountDescription    = strings.Repeat("亜", 256)
		justCountName           = strings.Repeat("亜", 32)
		justCountResult         = strings.Repeat("亜", 32)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		tooLongResultKanji      = strings.Repeat("亜", 33)
		invalidURL              = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		reqbody    schema.AddContestTeamRequest
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: description,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			schema.ContestTeam{
				Id:     dummyUUID(t), //テスト時にOptSyncIDで同期するため適当
				Name:   name,
				Result: result,
			},
		},
		"201 with kanji": {
			http.StatusCreated,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: justCountDescription,
				Link:        &link,
				Name:        justCountName,
				Result:      &justCountResult,
			},
			schema.ContestTeam{
				Id:     dummyUUID(t),
				Name:   justCountName,
				Result: justCountResult,
			},
		},
		"400 invalid description": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: tooLongString,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid description kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: tooLongDescriptionKanji,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: description,
				Link:        &invalidURL,
				Name:        name,
				Result:      &result,
			},
			httpError(t, "Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: description,
				Link:        &link,
				Name:        tooLongString,
				Result:      &result,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Name kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: description,
				Link:        &link,
				Name:        tooLongName,
				Result:      &result,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Result": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			schema.AddContestTeamRequest{
				Description: description,
				Link:        &link,
				Name:        name,
				Result:      &tooLongResultKanji,
			},
			httpError(t, "Bad Request: validate error: result: the length must be no more than 32."),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			schema.AddContestTeamRequest{
				Description: description,
				Link:        &link,
				Name:        name,
				Result:      &result,
			},
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodPost, e.URL(api.Contest.AddContestTeam, tt.contestID), &tt.reqbody)
			switch tt.want.(type) {
			case schema.ContestDetail:
				assertResponse(t, tt.statusCode, tt.want, res, optSyncID)
			case error:
				assertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// EditContestTeam PATCH /contests/:contestID/teams/:teamID
func TestEditContestTeam(t *testing.T) {
	var (
		description             = random.AlphaNumeric()
		link                    = random.RandURLString()
		name                    = random.AlphaNumeric()
		result                  = random.AlphaNumeric()
		tooLongString           = strings.Repeat("a", 260)
		justCountDescription    = strings.Repeat("亜", 256)
		justCountName           = strings.Repeat("亜", 32)
		justCountResult         = strings.Repeat("亜", 32)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		tooLongResultKanji      = strings.Repeat("亜", 33)
		invalidURL              = "invalid url"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		teamID     uuid.UUID
		reqBody    schema.EditContestTeamRequest
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Description: &description,
				Link:        &link,
				Name:        &name,
				Result:      &result,
			},
			nil,
		},
		"204 with kanji": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID2(),
			schema.EditContestTeamRequest{
				Description: &justCountDescription,
				Link:        &link,
				Name:        &justCountName,
				Result:      &justCountResult,
			},
			nil,
		},
		"204 without change": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID3(),
			schema.EditContestTeamRequest{},
			nil,
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid contestTeamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			schema.EditContestTeamRequest{},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid description": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Description: &tooLongString,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid description with kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Description: &tooLongDescriptionKanji,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 invalid Link": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Link: &invalidURL,
			},
			httpError(t, "Bad Request: validate error: link: must be a valid URL."),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Name: &tooLongString,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Name with kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Name: &tooLongName,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Result": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Result: &tooLongString,
			},
			httpError(t, "Bad Request: validate error: result: the length must be no more than 32."),
		},
		"400 invalid Result with kanji": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.EditContestTeamRequest{
				Result: &tooLongResultKanji,
			},
			httpError(t, "Bad Request: validate error: result: the length must be no more than 32."),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			schema.EditContestTeamRequest{
				Description: &description,
				Link:        &link,
				Name:        &name,
				Result:      &result,
			},
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var contestTeam schema.ContestTeamDetail
				res := doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeam, tt.contestID, tt.teamID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &contestTeam)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない

				// Update & Assert
				res = doRequest(t, e, http.MethodPatch, e.URL(api.Contest.EditContestTeam, tt.contestID, tt.teamID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)

				// Get updated response & Assert
				if tt.reqBody.Description != nil {
					contestTeam.Description = *tt.reqBody.Description
				}
				if tt.reqBody.Link != nil {
					contestTeam.Link = *tt.reqBody.Link
				}
				if tt.reqBody.Name != nil {
					contestTeam.Name = *tt.reqBody.Name
				}
				if tt.reqBody.Result != nil {
					contestTeam.Result = *tt.reqBody.Result
				}
				res = doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeam, tt.contestID, tt.teamID), nil)
				assertResponse(t, http.StatusOK, contestTeam, res)
			} else {
				res := doRequest(t, e, http.MethodPatch, e.URL(api.Contest.EditContestTeam, tt.contestID, tt.teamID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// DeleteContestTeam DELETE /contests/:contestID/teams/:teamID
func TestDeleteContestTeam(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		contestID  uuid.UUID
		teamID     uuid.UUID
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			nil,
		},
		"400: invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.ContestTeamID1(),
			httpError(t, "Bad Request: nil id"),
		},
		"400: invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
		"404: contest not found": {
			http.StatusNotFound,
			random.UUID(),
			mockdata.ContestTeamID1(),
			httpError(t, "Not Found: not found"),
		},
		"404: team not found": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodDelete, e.URL(api.Contest.DeleteContestTeam, tt.contestID, tt.teamID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
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
			[]schema.User{
				mockdata.CloneHandlerMockUsers()[0],
			},
		},
		"200 with no members": {
			http.StatusOK,
			mockdata.ContestID1(),
			mockdata.ContestTeamID2(),
			[]schema.User{},
		},
		"400 invalid contestID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.ContestTeamID1(),
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
		"404 contestID not exist": {
			http.StatusNotFound,
			random.UUID(),
			mockdata.ContestTeamID1(),
			httpError(t, "Not Found: not found"),
		},
		"404 teamID not exist": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeamMembers, tt.contestID, tt.teamID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
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
		reqbody    schema.MemberIDs
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.MemberIDs{
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
			schema.MemberIDs{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			schema.MemberIDs{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid memberID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.MemberIDs{
				Members: []uuid.UUID{
					uuid.Nil,
				},
			},
			httpError(t, "Bad Request: validate error: members: (0: must be a valid UUID v4.)."),
		},
		"400 invalid member": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.MemberIDs{
				Members: []uuid.UUID{
					random.UUID(),
				},
			},
			httpError(t, "Bad Request: argument error"),
		},
		"404 team not found": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			schema.MemberIDs{
				Members: []uuid.UUID{
					mockdata.UserID2(),
				},
			},
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodPost, e.URL(api.Contest.AddContestTeamMembers, tt.contestID, tt.teamID), &tt.reqbody)
			assertResponse(t, tt.statusCode, tt.want, res)
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
		reqbody    schema.MemberIDs
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.MemberIDs{
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
			schema.MemberIDs{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid teamID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			uuid.Nil,
			schema.MemberIDs{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid memberID": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.MemberIDs{
				Members: []uuid.UUID{
					uuid.Nil,
				},
			},
			httpError(t, "Bad Request: validate error: members: (0: must be a valid UUID v4.)."),
		},
		"400 invalid member": {
			http.StatusBadRequest,
			mockdata.ContestID1(),
			mockdata.ContestTeamID1(),
			schema.MemberIDs{
				Members: []uuid.UUID{
					random.UUID(),
				},
			},
			httpError(t, "Bad Request: argument error"),
		},
		"404 team not found": {
			http.StatusNotFound,
			mockdata.ContestID1(),
			random.UUID(),
			schema.MemberIDs{
				Members: []uuid.UUID{
					mockdata.UserID1(),
					mockdata.UserID2(),
				},
			},
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Update & Assert
				res := doRequest(t, e, http.MethodPut, e.URL(api.Contest.EditContestTeamMembers, tt.contestID, tt.teamID), &tt.reqbody)
				assertResponse(t, tt.statusCode, tt.want, res)

				// Assert
				res = doRequest(t, e, http.MethodGet, e.URL(api.Contest.GetContestTeamMembers, tt.contestID, tt.teamID), nil)
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
				res := doRequest(t, e, http.MethodPut, e.URL(api.Contest.EditContestTeamMembers, tt.contestID, tt.teamID), &tt.reqbody)
				assertResponse(t, tt.statusCode, tt.want, res)
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
	api, err := setupRoutes(t, e,)
	assert.NoError(t, err)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
                        t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.XXX.GetXXX, tt.userID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
*/
