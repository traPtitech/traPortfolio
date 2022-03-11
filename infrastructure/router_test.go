////go:build integration && router

package infrastructure_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

var (
	port = func() int {
		if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			return p
		}

		return 1323
	}()
	baseURL = fmt.Sprintf("http://localhost:%d/api/v1/", port)
)

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
		infrastructure.Setup(e, api)

		log.Fatal(e.Start(fmt.Sprintf(":%d", port)))
	}(api)

	time.Sleep(time.Second)

	os.Exit(m.Run())
}

func Test_Ping(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       []byte
	}{
		"200": {200, []byte("pong")},
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
