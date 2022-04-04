package testutils

import (
	"encoding/json"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func AssertResBody(t *testing.T, expected interface{}, res *testutil.CompletedRequest) bool {
	var expectedStr string
	switch v := expected.(type) {
	case string:
		expectedStr = v
	case []byte:
		expectedStr = string(v)
	default:
		expectedStr = string(mustMarshal(t, expected))
	}

	return assert.JSONEq(t, expectedStr, res.Recorder.Body.String())
}

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	return b
}
