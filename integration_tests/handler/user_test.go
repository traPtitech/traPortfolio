package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

// GetUsers GET /users
func TestGetUsers(t *testing.T) {
	var (
		includeSuspended = schema.IncludeSuspendedInQuery(true)
		name             = schema.NameInQuery(mockdata.MockUsers[0].Name)
		limitBlank       = schema.LimitInQuery(0)
		limitLessThan1   = schema.LimitInQuery(-1)
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqBody    schema.GetUsersParams
		want       interface{} // []schema.User | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			schema.GetUsersParams{},
			[]schema.User{
				mockdata.HMockUsers[0],
				mockdata.HMockUsers[2],
			},
		},
		"200 with includeSuspended": {
			http.StatusOK,
			schema.GetUsersParams{
				IncludeSuspended: &includeSuspended,
			},
			[]schema.User{
				mockdata.HMockUsers[0],
				mockdata.HMockUsers[1],
				mockdata.HMockUsers[2],
			},
		},
		"200 with name": {
			http.StatusOK,
			schema.GetUsersParams{
				Name: &name,
			},
			[]schema.User{
				mockdata.HMockUsers[0],
			},
		},
		"400 multiple params": {
			http.StatusBadRequest,
			schema.GetUsersParams{
				IncludeSuspended: &includeSuspended,
				Name:             &name,
			},
			testutils.HTTPError(t, "Bad Request: validate error: include_suspended and name cannot be specified at the same time"),
		},
		"400 invalid limit with 0": {
			http.StatusBadRequest,
			schema.GetUsersParams{
				Limit: &limitBlank,
			},
			testutils.HTTPError(t, "Bad Request: validate error: limit: cannot be blank."),
		},
		"400 invalid limit less than 1": {
			http.StatusBadRequest,
			schema.GetUsersParams{
				Limit: &limitLessThan1,
			},
			testutils.HTTPError(t, "Bad Request: validate error: limit: must be no less than 1."),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_users")
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
		want       interface{} // schema.UserDetail | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.UserID1(),
			mockdata.HMockUserDetails[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user")
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
		reqBody    schema.EditUserJSONRequestBody
		want       interface{} // nil or error
	}{
		"204": {
			http.StatusNoContent,
			mockdata.UserID1(),
			schema.EditUserJSONRequestBody{
				Bio:   &bio,
				Check: &check,
			},
			nil,
		},
		"204 without changes": {
			http.StatusNoContent,
			mockdata.UserID2(),
			schema.EditUserJSONRequestBody{},
			nil,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			schema.EditUserJSONRequestBody{},
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_update_user")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var user schema.UserDetail
				res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUser, tt.userID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &user)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない
				// Update & Assert
				res = testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.UpdateUser, tt.userID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
				// Get updated response & Assert
				if tt.reqBody.Check != nil && *tt.reqBody.Check == false {
					user.RealName = ""
				}
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
			mockdata.UserID1(),
			mockdata.HMockUserAccountsByID[mockdata.UserID1()],
		},
		"200 no accounts with existing userID": {
			http.StatusOK,
			mockdata.UserID2(),
			[]schema.Account{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user_accounts")
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
			mockdata.UserID1(),
			mockdata.AccountID1(),
			mockdata.HMockUserAccountsByID[mockdata.UserID1()][0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.AccountID1(),
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"400 invalid accountID": {
			http.StatusBadRequest,
			mockdata.UserID1(),
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"404 userID not found": {
			http.StatusNotFound,
			random.UUID(),
			mockdata.AccountID1(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
		"404 accountID not found": {
			http.StatusNotFound,
			mockdata.UserID1(),
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
		"404 both userID and accountID not found": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user_account")
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
		displayName          = random.AlphaNumeric()
		justCountDisplayName = strings.Repeat("亜", 256)
		tooLongDisplayName   = strings.Repeat("亜", 257)
		prPermitted          = schema.PrPermitted(random.Bool())
		testUserID           = mockdata.UserID1()
		accountType          = schema.AccountType(mockdata.AccountTypesMockUserDoesntHave(testUserID)[0])
		accountURL           = random.AccountURLString(domain.AccountType(accountType))
		conflictType         = schema.AccountType(mockdata.AccountTypesMockUserHas(testUserID)[0])
		testUserID2          = mockdata.UserID2()
		accountType2         = schema.AccountType(mockdata.AccountTypesMockUserDoesntHave(testUserID2)[0])
		accountURL2          = random.AccountURLString(domain.AccountType(accountType2))
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		reqBody    schema.AddUserAccountJSONRequestBody
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			testUserID,
			schema.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: prPermitted,
				Type:        accountType,
				Url:         accountURL,
			},
			schema.Account{
				Id:          uuid.Nil,
				DisplayName: displayName,
				PrPermitted: prPermitted,
				Type:        accountType,
				Url:         accountURL,
			},
		},
		"201 with kanji": {
			http.StatusCreated,
			testUserID2,
			schema.AddUserAccountJSONRequestBody{
				DisplayName: justCountDisplayName,
				PrPermitted: prPermitted,
				Type:        accountType2,
				Url:         accountURL2,
			},
			schema.Account{
				Id:          uuid.Nil,
				DisplayName: justCountDisplayName,
				PrPermitted: prPermitted,
				Type:        accountType2,
				Url:         accountURL2,
			},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			schema.AddUserAccountJSONRequestBody{},
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"400 invalid URL": {
			http.StatusBadRequest,
			testUserID,
			schema.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: prPermitted,
				Type:        accountType,
				Url:         "invalid url",
			},
			testutils.HTTPError(t, "Bad Request: validate error: url: must be a valid URL."),
		},
		"400 invalid account type": {
			http.StatusBadRequest,
			testUserID,
			schema.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: prPermitted,
				Type:        schema.AccountType(domain.AccountLimit),
				Url:         accountURL,
			},
			testutils.HTTPError(t, "Bad Request: validate error: type: must be no greater than 11."),
		},
		"409 conflict already exists": {
			http.StatusConflict,
			testUserID,
			schema.AddUserAccountJSONRequestBody{
				DisplayName: displayName,
				PrPermitted: prPermitted,
				Type:        conflictType,
				Url:         random.AccountURLString(domain.AccountType(conflictType)),
			},
			testutils.HTTPError(t, "Conflict: already exists"),
		},
		"400 too long display name": {
			http.StatusBadRequest,
			testUserID,
			schema.AddUserAccountJSONRequestBody{
				DisplayName: tooLongDisplayName,
				PrPermitted: prPermitted,
				Type:        accountType,
				Url:         accountURL,
			},
			testutils.HTTPError(t, "Bad Request: validate error: displayName: the length must be between 1 and 256."),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_add_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), &tt.reqBody)
			switch want := tt.want.(type) {
			case schema.Account:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res, testutils.OptSyncID, testutils.OptRetrieveID(&want.Id))
			case error:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// EditUserAccount PATCH /users/:userID/accounts/:accountID
func TestEditUserAccount(t *testing.T) {
	var (
		displayName        = random.AlphaNumeric()
		prPermitted        = schema.PrPermitted(random.Bool())
		testAccount        = mockdata.UserID1()
		accountType        = schema.AccountType(mockdata.AccountTypesMockUserHas(testAccount)[0])
		accountURL         = random.AccountURLString(domain.AccountType(accountType))
		initialAccountType = domain.AccountType(mockdata.AccountTypesMockUserDoesntHave(testAccount)[0])
		invalidAccountType = schema.AccountType(domain.GITHUB)
		invalidAccountURL  = random.RandURLString()
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode    int
		userID        uuid.UUID
		accountID     uuid.UUID
		reqBody       schema.EditUserAccountJSONRequestBody
		want          interface{} // nil | error
		needInsertion bool
	}{
		"204": {
			http.StatusNoContent,
			mockdata.UserID1(),
			mockdata.AccountID1(),
			schema.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
				PrPermitted: &prPermitted,
				Type:        &accountType,
				Url:         &accountURL,
			},
			nil,
			false,
		},
		"204 without changes": { // TODO: https://github.com/traPtitech/traPortfolio/issues/292
			http.StatusNoContent,
			mockdata.UserID2(),
			random.UUID(),
			schema.EditUserAccountJSONRequestBody{},
			nil,
			true,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			mockdata.AccountID1(),
			schema.EditUserAccountJSONRequestBody{},
			testutils.HTTPError(t, "Bad Request: nil id"),
			false,
		},
		"400 invalid accountID": {
			http.StatusBadRequest,
			mockdata.UserID1(),
			uuid.Nil,
			schema.EditUserAccountJSONRequestBody{},
			testutils.HTTPError(t, "Bad Request: nil id"),
			false,
		},
		"400 invalud url without accountType": {
			http.StatusBadRequest,
			mockdata.UserID1(),
			mockdata.AccountID1(),
			schema.EditUserAccountJSONRequestBody{
				Url: &invalidAccountURL,
			},
			testutils.HTTPError(t, "Bad Request: argument error"),
			false,
		},
		"400 invalid url without accountURL": {
			http.StatusBadRequest,
			mockdata.UserID1(),
			mockdata.AccountID1(),
			schema.EditUserAccountJSONRequestBody{
				Type: &invalidAccountType,
			},
			testutils.HTTPError(t, "Bad Request: argument error"),
			false,
		},
		"404 user not found": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			schema.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
			},
			testutils.HTTPError(t, "Not Found: not found"),
			false,
		},
		"404 account not found": {
			http.StatusNotFound,
			mockdata.UserID1(),
			random.UUID(),
			schema.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
			},
			testutils.HTTPError(t, "Not Found: not found"),
			false,
		},
		"404 account type conflicted by update": {
			http.StatusConflict,
			mockdata.UserID1(),
			mockdata.AccountID1(),
			schema.EditUserAccountJSONRequestBody{
				DisplayName: &displayName,
				PrPermitted: &prPermitted,
				Type:        &accountType,
				Url:         &accountURL,
			},
			testutils.HTTPError(t, "Conflict: already exists"),
			true,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_edit_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			account := schema.Account{}
			if tt.needInsertion {
				// Insert & Assert
				account = schema.Account{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: schema.PrPermitted(random.Bool()),
					Type:        schema.AccountType(initialAccountType),
					Url:         random.AccountURLString(initialAccountType),
				}
				res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), schema.AddUserAccountJSONRequestBody{
					DisplayName: account.DisplayName,
					PrPermitted: account.PrPermitted,
					Type:        account.Type,
					Url:         account.Url,
				})
				testutils.AssertResponse(t, http.StatusCreated, account, res, testutils.OptSyncID, testutils.OptRetrieveID(&tt.accountID))
				account.Id = tt.accountID
			} else {
				// Get account data
				res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetUserAccount, tt.userID, tt.accountID), nil)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &account))
			}
			// Update & Assert
			res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.User.EditUserAccount, tt.userID, tt.accountID), tt.reqBody)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			if tt.statusCode == http.StatusNoContent {
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
			mockdata.UserID1(),
			testutils.DummyUUID(t),
			nil,
			true,
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			random.UUID(),
			testutils.HTTPError(t, "Bad Request: nil id"),
			false,
		},
		"404 user not found": {
			http.StatusNotFound,
			random.UUID(),
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
			false,
		},
		"404 account not found": {
			http.StatusNotFound,
			mockdata.UserID1(),
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
			false,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_delete_user_account")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.needInsertion {
				accountType := mockdata.AccountTypesMockUserDoesntHave(tt.userID)[0]
				reqBody := schema.AddUserAccountJSONRequestBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: schema.PrPermitted(random.Bool()),
					Type:        accountType,
					Url:         random.AccountURLString(domain.AccountType(accountType)),
				}
				res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.User.AddUserAccount, tt.userID), &reqBody)
				testutils.AssertResponse(t, http.StatusCreated, schema.Account{
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
			mockdata.UserID1(),
			[]schema.UserProject{mockdata.HMockUserProjects[0]},
		},
		"200 no projects with existing userID": {
			http.StatusOK,
			mockdata.UserID3(),
			[]schema.Project{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user_projects")
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
			mockdata.UserID1(),
			mockdata.HMockUserContestsByID[mockdata.UserID1()],
		},
		"200 no contests with existing userID": {
			http.StatusOK,
			mockdata.UserID2(),
			[]schema.Contest{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user_contests")
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
			mockdata.UserID1(),
			mockdata.HMockUserGroupsByID[mockdata.UserID1()],
		},
		"200 no groups with existing userID": {
			http.StatusOK,
			mockdata.UserID2(),
			[]schema.Group{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
		"404 no accounts with not-existing userID": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user_groups")
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
			mockdata.UserID1(),
			mockdata.HMockUserEvents[:1],
		},
		"200 no events with existing userID": {
			http.StatusOK,
			mockdata.UserID3(),
			[]schema.Event{},
		},
		"200 no events with non-existing userID": {
			http.StatusOK,
			random.UUID(),
			[]schema.Event{},
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError(t, "Bad Request: nil id"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName(t, "user_handler_get_user_events")
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
