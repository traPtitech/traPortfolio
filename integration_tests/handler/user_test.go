//go:build integration && db

package handler

import (
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
		includeSuspended handler.IncludeSuspendedInQuery = true
		name             handler.NameInQuery             = handler.NameInQuery(mockdata.MockUsers[0].Name)
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
				mockdata.HMockUser1,
				mockdata.HMockUser3,
			},
		},
		"200 with includeSuspended": {
			http.StatusOK,
			handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
			},
			[]handler.User{
				mockdata.HMockUser1,
				mockdata.HMockUser2,
				mockdata.HMockUser3,
			},
		},
		"200 with name": {
			http.StatusOK,
			handler.GetUsersParams{
				Name: &name,
			},
			[]handler.User{
				mockdata.HMockUser1,
			},
		},
		"400 multiple params": {
			http.StatusBadRequest,
			handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
				Name:             &name,
			},
			handler.ConvertError(t, repository.ErrInvalidArg),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_users")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
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
			mockdata.HMockUserDetail1.Id,
			mockdata.HMockUserDetail1,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrValidate),
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
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUser, tt.userID.String()), nil)
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
		reqBody    handler.EditUser
		want       interface{} // nil or error
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockUser1.Id,
			handler.EditUser{
				Bio:   &bio,
				Check: &check,
			},
			nil,
		},
		"204 without changes": {
			http.StatusNoContent,
			mockdata.HMockUser1.Id,
			handler.EditUser{},
			nil,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.EditUser{},
			handler.ConvertError(t, repository.ErrValidate),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_update_user")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.UpdateUser, tt.userID.String()), &tt.reqBody)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
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
			mockdata.HMockUser1.Id,
			[]handler.Account{
				mockdata.HMockAccount,
			},
		},
		"200 no accounts": {
			http.StatusOK,
			random.UUID(),
			[]handler.Account{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrValidate),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_get_user_accounts")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
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
			mockdata.HMockUser1.Id,
			mockdata.HMockAccount.Id,
			mockdata.HMockAccount,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.HMockAccount.Id,
			handler.ConvertError(t, repository.ErrValidate),
		},
		"400 invalid accountID": {
			http.StatusBadRequest,
			mockdata.HMockUser1.Id,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrValidate),
		},
		"404": {
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
			mockdata.HMockUser1.Id,
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
			handler.ConvertError(t, repository.ErrValidate),
		},
		"400 invalid URL": {
			http.StatusBadRequest,
			mockdata.HMockUser1.Id,
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
			mockdata.HMockUser1.Id,
			handler.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: handler.PrPermitted(prPermitted),
				Type:        handler.AccountType(domain.AccountLimit),
				Url:         url,
			},
			handler.ConvertError(t, repository.ErrInvalidArg),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("user_handler_add_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), &tt.reqBody)
			switch tt.want.(type) {
			case handler.Account:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res, testutils.OptSyncID)
			case error:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// // EditUserAccount PATCH /users/:userID/accounts/:accountID
// func TestEditUserAccount(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		statusCode int
// 		userID     uuid.UUID
// 		accountID  uuid.UUID
// 		reqBody    handler.EditUserAccountJSONRequestBody
// 		want       interface{}
// 	}{
// 		// TODO: Add cases
// 	}

// 	e := echo.New()
// 	conf := testutils.GetConfigWithDBName("user_handler_edit_user_account")
// 	api, err := testutils.SetupRoutes(t, e, conf)
// 	assert.NoError(t, err)
// 	for name, tt := range tests {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.EditUserAccount, tt.userID, tt.accountID), tt.reqBody)
// 			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
// 		})
// 	}
// }

// // DeleteUserAccount DELETE /users/:userID/accounts/:accountID
// func TestDeleteUserAccount(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		statusCode int
// 		userID     uuid.UUID
// 		accountID  uuid.UUID
// 		want       interface{}
// 	}{
// 		// TODO: Add cases
// 	}

// 	e := echo.New()
// 	conf := testutils.GetConfigWithDBName("user_handler_delete_user_account")
// 	api, err := testutils.SetupRoutes(t, e, conf)
// 	assert.NoError(t, err)
// 	for name, tt := range tests {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			res := testutils.DoRequest(t, e, http.MethodDelete, e.URL(api.User.DeleteUserAccount, tt.userID, tt.accountID), nil)
// 			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
// 		})
// 	}
// }

// // GetUserProjects GET /users/:userID/projects
// func TestGetUserProjects(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		statusCode int
// 		userID     uuid.UUID
// 		want       interface{}
// 	}{
// 		// TODO: Add cases
// 	}

// 	e := echo.New()
// 	conf := testutils.GetConfigWithDBName("user_handler_get_user_projects")
// 	api, err := testutils.SetupRoutes(t, e, conf)
// 	assert.NoError(t, err)
// 	for name, tt := range tests {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserProjects, tt.userID), nil)
// 			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
// 		})
// 	}
// }

// // GetUserContests GET /users/:userID/contests
// func TestGetUserContests(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		statusCode int
// 		userID     uuid.UUID
// 		want       interface{}
// 	}{
// 		// TODO: Add cases
// 	}

// 	e := echo.New()
// 	conf := testutils.GetConfigWithDBName("user_handler_get_user_contests")
// 	api, err := testutils.SetupRoutes(t, e, conf)
// 	assert.NoError(t, err)
// 	for name, tt := range tests {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserContests, tt.userID), nil)
// 			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
// 		})
// 	}
// }

// // GetUserGroups GET /users/:userID/groups
// func TestGetUserGroups(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		statusCode int
// 		userID     uuid.UUID
// 		want       interface{}
// 	}{
// 		// TODO: Add cases
// 	}

// 	e := echo.New()
// 	conf := testutils.GetConfigWithDBName("user_handler_get_user_groups")
// 	api, err := testutils.SetupRoutes(t, e, conf)
// 	assert.NoError(t, err)
// 	for name, tt := range tests {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserGroups, tt.userID), nil)
// 			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
// 		})
// 	}
// }

// // GetUserEvents GET /users/:userID/events
// func TestGetUserEvents(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		statusCode int
// 		userID     uuid.UUID
// 		want       interface{}
// 	}{
// 		// TODO: Add cases
// 	}

// 	e := echo.New()
// 	conf := testutils.GetConfigWithDBName("user_handler_get_user_events")
// 	api, err := testutils.SetupRoutes(t, e, conf)
// 	assert.NoError(t, err)
// 	for name, tt := range tests {
// 		tt := tt
// 		t.Run(name, func(t *testing.T) {
// 			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserEvents, tt.userID), nil)
// 			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
// 		})
// 	}
// }
