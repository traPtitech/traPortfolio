//go:build integration && db

package handler_test

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
)

var (
	sampleUser1 = handler.User{
		Id:       mockdata.MockUsers[0].ID,
		Name:     mockdata.MockUsers[0].Name,
		RealName: mockdata.MockPortalUsers[0].RealName,
	}
	sampleUser2 = handler.User{
		Id:       mockdata.MockUsers[1].ID,
		Name:     mockdata.MockUsers[1].Name,
		RealName: mockdata.MockPortalUsers[1].RealName,
	}
	sampleUser3 = handler.User{
		Id:       mockdata.MockUsers[2].ID,
		Name:     mockdata.MockUsers[2].Name,
		RealName: mockdata.MockPortalUsers[2].RealName,
	}

	sampleAccount = handler.Account{
		DisplayName: mockdata.MockAccount.Name,
		Id:          mockdata.MockAccount.ID,
		PrPermitted: handler.PrPermitted(mockdata.MockAccount.Check),
		Type:        handler.AccountType(mockdata.MockAccount.Type),
		Url:         mockdata.MockAccount.URL,
	}

	sampleUserDetail1 = handler.UserDetail{
		User:     sampleUser1,
		Accounts: []handler.Account{sampleAccount},
		Bio:      mockdata.MockUsers[0].Description,
		State:    handler.UserAccountState(mockdata.MockTraQUsers[0].User.State),
	}
)

// GET /users
func TestGetUsers(t *testing.T) {
	var (
		includeSuspended handler.IncludeSuspendedInQuery = true
		name             handler.NameInQuery             = handler.NameInQuery(mockdata.MockUsers[0].Name)
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		params     handler.GetUsersParams
		want       interface{} // []handler.User | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			handler.GetUsersParams{},
			[]handler.User{
				sampleUser1,
				sampleUser3,
			},
		},
		"200 with includeSuspended": {
			http.StatusOK,
			handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
			},
			[]handler.User{
				sampleUser1,
				sampleUser2,
				sampleUser3,
			},
		},
		"200 with name": {
			http.StatusOK,
			handler.GetUsersParams{
				Name: &name,
			},
			[]handler.User{
				sampleUser1,
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
	conf := testutils.GetConfigWithDBName("get_users")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetAll), &tt.params)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GET /users/:userID
func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		userID     uuid.UUID
		want       interface{} // handler.UserDetail | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			sampleUserDetail1.Id,
			sampleUserDetail1,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("get_user")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.User.GetByID, tt.userID.String()), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
