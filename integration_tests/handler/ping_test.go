package handler

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
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
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Ping.Ping), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
