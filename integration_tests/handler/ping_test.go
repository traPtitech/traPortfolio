//go:build integration && db

package handler_test

import (
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
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
	api, err := testutils.SetupRoutes(t, e, "get_ping", nil)
	assert.NoError(t, err)
	req := testutil.NewRequest().Get(e.URL(api.Ping.Ping))
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := req.Go(t, e)
			assert.Equal(t, tt.statusCode, res.Code())
			assert.Equal(t, tt.want, res.Recorder.Body.Bytes())
		})
	}
}
