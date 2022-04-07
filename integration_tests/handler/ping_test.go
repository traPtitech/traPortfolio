//go:build integration && db

package handler_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
)

// GET /ping
func TestPing(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       []byte
	}{
		"200": {http.StatusOK, []byte("pong")},
	}

	e := echo.New()

	conf := testutils.GetConfigWithDBName("get_ping")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Ping.Ping), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
