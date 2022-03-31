//go:build integration && db

package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cli, err := NewClientWithResponses(baseURL)
			assert.NoError(t, err)

			res, err := cli.PingWithResponse(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, tt.statusCode, res.StatusCode())
			assert.Equal(t, tt.want, res.Body)
		})
	}
}
