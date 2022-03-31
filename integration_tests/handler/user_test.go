//go:build integration && db

package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

var (
	sampleUser1 = User{
		Id:       uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
		Name:     "user1",
		RealName: "ユーザー1 ユーザー1",
	}
	sampleUser2 = User{
		Id:       uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
		Name:     "user2",
		RealName: "ユーザー2 ユーザー2",
	}
	sampleUser3 = User{
		Id:       uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
		Name:     "lolico",
		RealName: "東 工子",
	}
)

// GET /users
func TestGetUsers(t *testing.T) {
	var (
		includeSuspended IncludeSuspendedInQuery = true
		name             NameInQuery             = "user1"
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		params     GetUsersParams
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			GetUsersParams{},
			[]User{
				sampleUser1,
				sampleUser3,
			},
		},
		"200 with includeSuspended": {
			http.StatusOK,
			GetUsersParams{
				IncludeSuspended: &includeSuspended,
			},
			[]User{
				sampleUser1,
				sampleUser2,
				sampleUser3,
			},
		},
		"200 with name": {
			http.StatusOK,
			GetUsersParams{
				Name: &name,
			},
			[]User{
				sampleUser1,
			},
		},
		"400 multiple params": {
			http.StatusBadRequest,
			GetUsersParams{
				IncludeSuspended: &includeSuspended,
				Name:             &name,
			},
			handler.ConvertError(t, repository.ErrInvalidArg),
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cli, err := NewClientWithResponses(baseURL)
			assert.NoError(t, err)

			res, err := cli.GetUsersWithResponse(context.Background(), &tt.params)
			assert.NoError(t, err)

			assert.Equal(t, tt.statusCode, res.StatusCode())
			assert.JSONEq(t, string(mustMarshal(tt.want)), string(res.Body))
		})
	}
}
