//go:build integration && db

package testutils

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func AssertResponse(t *testing.T, expectedStatusCode int, expectedBody interface{}, res *httptest.ResponseRecorder, opts ...Option) {
	t.Helper()

	assert.Equal(t, expectedStatusCode, res.Code)

	for _, o := range opts {
		assert.NoError(t, o(t, &expectedBody, res))
	}

	actual := res.Body
	switch expected := expectedBody.(type) {
	case string:
		assert.Equal(t, expected, actual.String())
	case []byte:
		assert.Equal(t, expected, actual.Bytes())
	case nil:
		assert.Empty(t, actual.String())
	default:
		b, err := json.Marshal(expected)
		assert.NoError(t, err)
		assert.JSONEq(t, string(b), actual.String())
	}
}

// NOTE: expectedBodyPtr must be a pointer to expectedBody
type Option func(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error

// TODO: testifyからgo-cmpに乗り換えたらFilterValuesを使ってIDをignoreする
func OptSyncID(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error {
	t.Helper()

	m := struct {
		ID uuid.UUID `json:"id"`
	}{}

	if err := json.Unmarshal(res.Body.Bytes(), &m); err != nil {
		return err
	}

	v := reflect.ValueOf(expectedBodyPtr).Elem()
	tmp := reflect.New(v.Elem().Type()).Elem()
	tmp.Set(v.Elem())
	tmp.FieldByName("Id").Set(reflect.ValueOf(m.ID))
	v.Set(tmp)

	return nil
}

func OptRetrieveID(idPtr *uuid.UUID) Option {
	return func(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error {
		t.Helper()

		m := struct {
			ID uuid.UUID `json:"id"`
		}{}

		if err := json.Unmarshal(res.Body.Bytes(), &m); err != nil {
			return err
		}

		*idPtr = m.ID

		return nil
	}
}
