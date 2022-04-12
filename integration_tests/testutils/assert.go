//go:build integration && db

package testutils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertResponse(t *testing.T, expectedStatusCode int, expectedBody interface{}, res *httptest.ResponseRecorder) {
	t.Helper()

	assert.Equal(t, expectedStatusCode, res.Code)

	actual := res.Body
	switch expected := expectedBody.(type) {
	case string:
		assert.Equal(t, expected, actual.String())
	case []byte:
		assert.Equal(t, expected, actual.Bytes())
	default:
		b, err := json.Marshal(expected)
		assert.NoError(t, err)
		assert.JSONEq(t, string(b), actual.String())
	}
}
