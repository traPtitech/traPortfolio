package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

// GetUsers GET /users
func TestGetUsers(t *testing.T) {
	var (
		includeSuspended = handler.IncludeSuspendedInQuery(true)
		name             = handler.NameInQuery(mockdata.MockUsers[0].Name)
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqBody    handler.GetUsersParams
		want       interface{} // []handler.User | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			handler.GetUsersParams{},
			[]handler.User{
				mockdata.HMockUsers[0],
				mockdata.HMockUsers[2],
			},
		},
		"200 with includeSuspended": {
			http.StatusOK,
			handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
			},
			[]handler.User{
				mockdata.HMockUsers[0],
				mockdata.HMockUsers[1],
				mockdata.HMockUsers[2],
			},
		},
		"200 with name": {
			http.StatusOK,
			handler.GetUsersParams{
				Name: &name,
			},
			[]handler.User{
				mockdata.HMockUsers[0],
			},
		},
		"400 multiple params": {
			http.StatusBadRequest,
			handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
				Name:             &name,
			},
			handler.ConvertError(t, repository.ErrValidate),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_users")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUsers), &tt.reqBody)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetUser GET /users/:userID
func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{} // handler.UserDetail | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUserDetails[0].Id,
			mockdata.HMockUserDetails[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUser, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// UpdateUser PATCH /users/:userID
func TestUpdateUser(t *testing.T) {
	var (
		bio   = random.AlphaNumeric()
		check = random.Bool()
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		reqBody    handler.EditUserRequest
		want       interface{} // nil or error
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockUsers[0].Id,
			handler.EditUserRequest{
				Bio:   &bio,
				Check: &check,
			},
			nil,
		},
		"204 without changes": {
			http.StatusNoContent,
			mockdata.HMockUsers[1].Id,
			handler.EditUserRequest{},
			nil,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.EditUserRequest{},
			handler.ConvertError(t, repository.ErrNilID),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_update_user")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var user handler.UserDetail
				res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUser, tt.userID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &user)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない
				// Update & Assert
				res = testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.UpdateUser, tt.userID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
				// Get updated response & Assert
				if tt.reqBody.Bio != nil {
					user.Bio = *tt.reqBody.Bio
				}
				// if tt.reqBody.Check != nil {} // TODO: Checkに応じて処理を書く
				res = testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUser, tt.userID), nil)
				testutils.AssertResponse(t, http.StatusOK, user, res)
			} else {
				res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.UpdateUser, tt.userID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// GetUserAccounts GET /users/:userID/accounts
func TestGetUserAccounts(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUsers[0].Id,
			mockdata.HMockUserAccounts,
		},
		"200 no accounts with existing userID": {
			http.StatusOK,
			mockdata.HMockUsers[1].Id,
			[]handler.Account{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_accounts")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserAccounts, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetUserAccount GET /users/:userID/accounts/:accountID
func TestGetUserAccount(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		accountID  uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUsers[0].Id,
			mockdata.HMockUserAccounts[0].Id,
			mockdata.HMockUserAccounts[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.HMockUserAccounts[0].Id,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"400 invalid accountID": {
			http.StatusBadRequest,
			mockdata.HMockUsers[0].Id,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404 userID not found": {
			http.StatusNotFound,
			random.UUID(),
			mockdata.HMockUserAccounts[0].Id,
			handler.ConvertError(t, repository.ErrNotFound),
		},
		"404 accountID not found": {
			http.StatusNotFound,
			mockdata.HMockUsers[0].Id,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
		"404 both userID and accountID not found": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserAccount, tt.userID, tt.accountID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// AddUserAccount POST /users/:userID/accounts
func TestAddUserAccount(t *testing.T) {
	var (
		displayName = random.AlphaNumeric()
		prPermitted = random.Bool()
		atype       = rand.Intn(int(domain.AccountLimit)) // TODO: openapiでenumを定義する
		url         = random.RandURLString()
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		reqBody    handler.AddUserAccountJSONRequestBody
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			mockdata.HMockUsers[0].Id,
			handler.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: handler.PrPermitted(prPermitted),
				Type:        handler.AccountType(atype),
				Url:         url,
			},
			handler.Account{
				Id:          uuid.Nil,
				DisplayName: displayName,
				PrPermitted: handler.PrPermitted(prPermitted),
				Type:        handler.AccountType(atype),
				Url:         url,
			},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.AddUserAccountJSONRequestBody{},
			handler.ConvertError(t, repository.ErrNilID),
		},
		"400 invalid URL": {
			http.StatusBadRequest,
			mockdata.HMockUsers[0].Id,
			handler.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: handler.PrPermitted(prPermitted),
				Type:        handler.AccountType(atype),
				Url:         "invalid url",
			},
			handler.ConvertError(t, repository.ErrValidate),
		},
		"400 invalid account type": {
			http.StatusBadRequest,
			mockdata.HMockUsers[0].Id,
			handler.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: handler.PrPermitted(prPermitted),
				Type:        handler.AccountType(domain.AccountLimit),
				Url:         url,
			},
			handler.ConvertError(t, repository.ErrValidate),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_add_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), &tt.reqBody)
			switch want := tt.want.(type) {
			case handler.Account:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res, testutils.OptSyncID, testutils.OptRetrieveID(&want.Id))
			case error:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// EditUserRequestAccount PATCH /users/:userID/accounts/:accountID
func TestEditUserRequestAccount(t *testing.T) {
	var (
		displayName = random.AlphaNumeric()
		prPermitted = random.Bool()
		atype       = int64(rand.Intn(int(domain.AccountLimit))) // TODO: openapiでenumを定義する
		url         = random.RandURLString()
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		accountID  uuid.UUID
		reqBody    handler.EditUserAccountJSONRequestBody
		want       interface{} // nil | error
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockUsers[0].Id,
			testutils.DummyUUID(),
			handler.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
				PrPermitted: (*handler.PrPermitted)(&prPermitted),
				Type:        (*handler.AccountType)(&atype),
				Url:         &url,
			},
			nil,
		},
		"204 without changes": { // TODO: https://github.com/traPtitech/traPortfolio/issues/292
			http.StatusNoContent,
			mockdata.HMockUsers[0].Id,
			testutils.DummyUUID(),
			handler.EditUserAccountJSONRequestBody{},
			nil,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.HMockUserAccounts[0].Id,
			handler.EditUserAccountJSONRequestBody{},
			handler.ConvertError(t, repository.ErrNilID),
		},
		"400 invalid accountID": {
			http.StatusBadRequest,
			mockdata.HMockUsers[0].Id,
			uuid.Nil,
			handler.EditUserAccountJSONRequestBody{},
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404 user not found": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			handler.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
			},
			handler.ConvertError(t, repository.ErrNotFound),
		},
		"404 account not found": {
			http.StatusNotFound,
			mockdata.HMockUsers[0].Id,
			random.UUID(),
			handler.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
			},
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_edit_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Insert & Assert
				account := handler.Account{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: handler.PrPermitted(random.Bool()),
					Type:        handler.AccountType(rand.Intn(int(domain.AccountLimit))),
					Url:         random.RandURLString(),
				}
				res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), handler.AddUserAccountJSONRequestBody{
					DisplayName: account.DisplayName,
					PrPermitted: account.PrPermitted,
					Type:        account.Type,
					Url:         account.Url,
				})
				testutils.AssertResponse(t, http.StatusCreated, account, res, testutils.OptSyncID, testutils.OptRetrieveID(&tt.accountID))
				account.Id = tt.accountID
				// Update & Assert
				res = testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.EditUserAccount, tt.userID, tt.accountID), tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
				// Get updated response & Assert
				if tt.reqBody.DisplayName != nil {
					account.DisplayName = *tt.reqBody.DisplayName
				}
				if tt.reqBody.PrPermitted != nil {
					account.PrPermitted = *tt.reqBody.PrPermitted
				}
				if tt.reqBody.Type != nil {
					account.Type = *tt.reqBody.Type
				}
				if tt.reqBody.Url != nil {
					account.Url = *tt.reqBody.Url
				}
				res = testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserAccount, tt.userID, tt.accountID), nil)
				testutils.AssertResponse(t, http.StatusOK, account, res)
			} else {
				res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.EditUserAccount, tt.userID, tt.accountID), tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// DeleteUserAccount DELETE /users/:userID/accounts/:accountID
func TestDeleteUserAccount(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode    int
		userID        uuid.UUID
		accountID     uuid.UUID
		want          interface{}
		needInsertion bool
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockUsers[0].Id,
			testutils.DummyUUID(),
			nil,
			true,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNilID),
			false,
		},
		"404 user not found": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
			false,
		},
		"404 account not found": {
			http.StatusNotFound,
			mockdata.HMockUsers[0].Id,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
			false,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_delete_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.needInsertion {
				reqBody := handler.AddUserAccountJSONRequestBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: handler.PrPermitted(random.Bool()),
					Type:        handler.AccountType(rand.Intn(int(domain.AccountLimit))),
					Url:         random.RandURLString(),
				}
				res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), &reqBody)
				testutils.AssertResponse(t, http.StatusCreated, handler.Account{
					DisplayName: reqBody.DisplayName,
					PrPermitted: reqBody.PrPermitted,
					Type:        reqBody.Type,
					Url:         reqBody.Url,
				}, res, testutils.OptSyncID, testutils.OptRetrieveID(&tt.accountID))
			}
			res := testutils.DoRequest(t, e, http.MethodDelete, e.URL(api.User.DeleteUserAccount, tt.userID, tt.accountID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetUserProjects GET /users/:userID/projects
func TestGetUserProjects(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUsers[0].Id,
			mockdata.HMockUserProjects,
		},
		"200 no projects with existing userID": {
			http.StatusOK,
			mockdata.HMockUsers[1].Id,
			[]handler.Project{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_projects")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserProjects, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetUserContests GET /users/:userID/contests
func TestGetUserContests(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUsers[0].Id,
			mockdata.HMockUserContests,
		},
		"200 no contests with existing userID": {
			http.StatusOK,
			mockdata.HMockUsers[1].Id,
			[]handler.Contest{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_contests")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserContests, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetUserGroups GET /users/:userID/groups
func TestGetUserGroups(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUsers[0].Id,
			mockdata.HMockUserGroups,
		},
		"200 no groups with existing userID": {
			http.StatusOK,
			mockdata.HMockUsers[1].Id,
			[]handler.Group{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_groups")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserGroups, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetUserEvents GET /users/:userID/events
func TestGetUserEvents(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockUsers[0].Id,
			mockdata.HMockUserEvents,
		},
		"200 no events with existing userID": {
			http.StatusOK,
			mockdata.HMockUsers[1].Id,
			[]handler.Event{
				mockdata.HMockUserEvents[1],
			},
		},
		"200 no events with non-existing userID": {
			http.StatusOK,
			random.UUID(),
			[]handler.Event{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrNilID),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_events")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserEvents, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
