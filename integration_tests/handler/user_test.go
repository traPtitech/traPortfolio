//go:build integration && db

package handler_test

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

var (
	sampleUser1 = handler.User{
		Id:       uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"), // TODO: 変数で管理する 以下も同様
		Name:     "user1",
		RealName: "ユーザー1 ユーザー1",
	}
	sampleUser2 = handler.User{
		Id:       uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
		Name:     "user2",
		RealName: "ユーザー2 ユーザー2",
	}
	sampleUser3 = handler.User{
		Id:       uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
		Name:     "lolico",
		RealName: "東 工子",
	}
)

func initUser(h database.SQLHandler) error {
	if err := h.Create([]*model.User{
		{
			ID:          sampleUser1.Id,
			Description: "I am user1",
			Check:       true,
			Name:        sampleUser1.Name,
		},
		{
			ID:          sampleUser2.Id,
			Description: "I am user2",
			Check:       true,
			Name:        sampleUser2.Name,
		},
		{
			ID:          sampleUser3.Id,
			Description: "I am lolico",
			Check:       false,
			Name:        sampleUser3.Name,
		},
	}).Error(); err != nil {
		return err
	}

	return nil
}

// GET /users
func TestGetUsers(t *testing.T) {
	var (
		includeSuspended handler.IncludeSuspendedInQuery = true
		name             handler.NameInQuery             = "user1"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		params     handler.GetUsersParams
		want       interface{}
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
	api, err := testutils.SetupRoutes(t, e, "get_users", initUser)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := testutils.DoRequest(t, e, http.MethodGet, api.User.GetAll, &tt.params)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
