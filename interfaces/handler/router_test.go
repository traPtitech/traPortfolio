//go:build integration && router

package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

var (
	port = func() int {
		if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			return p
		}

		return 1323
	}()
	baseURL = fmt.Sprintf("http://localhost:%d/api/v1/", port)

	sampleUser1 = handler.User{
		Id:       uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
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

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b
}

func TestMain(m *testing.M) {
	s := infrastructure.NewSQLConfig("root", "password", "localhost", "portfolio", 3307)
	t := infrastructure.NewTraQConfig("", "", true)
	p := infrastructure.NewPortalConfig("", "", true)
	k := infrastructure.NewKnoqConfig("", "", true)
	api, err := infrastructure.InjectAPIServer(&s, &t, &p, &k)
	if err != nil {
		log.Fatal(err)
	}

	go func(api handler.API) {
		e := echo.New()
		handler.Setup(e, api)

		log.Fatal(e.Start(fmt.Sprintf(":%d", port)))
	}(api)

	time.Sleep(time.Second)

	os.Exit(m.Run())
}

// GET /ping
func TestPing(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       []byte
	}{
		"200": {http.StatusOK, []byte("pong")},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cli, err := handler.NewClientWithResponses(baseURL)
			assert.NoError(t, err)

			res, err := cli.PingWithResponse(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, tt.statusCode, res.StatusCode())

			actual := reflect.ValueOf(*res).FieldByName("Body").Interface().([]byte)
			assert.Equal(t, tt.want, actual)
		})
	}
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
		params     *handler.GetUsersParams
		want       []byte
	}{
		"200": {
			http.StatusOK,
			new(handler.GetUsersParams),
			mustMarshal(&[]handler.User{
				sampleUser1,
				sampleUser3,
			}),
		},
		"200 with includeSuspended": {
			http.StatusOK,
			&handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
			},
			mustMarshal(&[]handler.User{
				sampleUser1,
				sampleUser2,
				sampleUser3,
			}),
		},
		"200 with name": {
			http.StatusOK,
			&handler.GetUsersParams{
				Name: &name,
			},
			mustMarshal(&[]handler.User{
				sampleUser1,
			}),
		},
		"400 multiple params": {
			http.StatusBadRequest,
			&handler.GetUsersParams{
				IncludeSuspended: &includeSuspended,
				Name:             &name,
			},
			mustMarshal(handler.ConvertError(repository.ErrInvalidArg)),
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cli, err := handler.NewClientWithResponses(baseURL)
			assert.NoError(t, err)

			res, err := cli.GetUsersWithResponse(context.Background(), tt.params)
			assert.NoError(t, err)

			assert.Equal(t, tt.statusCode, res.StatusCode())
			assert.JSONEq(t, string(tt.want), string(res.Body))
		})
	}
}
