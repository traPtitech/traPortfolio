package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

func setupRoutes(t *testing.T, e *echo.Echo, conf *config.Config) (*handler.API, error) {
	t.Helper()

	db := testutils.SetupGormDB(t, conf.SQLConf())
	if err := mockdata.InsertSampleDataToDB(db); err != nil {
		return nil, err
	}

	api, err := infrastructure.InjectAPIServer(conf, db)
	if err != nil {
		return nil, err
	}

	if err := handler.Setup(e, api); err != nil {
		return nil, err
	}

	return &api, nil
}

func doRequest(t *testing.T, e *echo.Echo, method string, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		assert.NoError(t, err)

		bodyReader = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}

// OptRetrieveIDなどですぐに変更され得るUUIDであることの明示に使う
func dummyUUID(t *testing.T) uuid.UUID {
	t.Helper()
	return random.UUID()
}

func httpError(t *testing.T, message string) *echo.HTTPError {
	t.Helper()
	return &echo.HTTPError{
		Message: message,
	}
}

func assertResponse(t *testing.T, expectedStatusCode int, expectedBody interface{}, res *httptest.ResponseRecorder, opts ...option) {
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
type option func(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error

// TODO: testifyからgo-cmpに乗り換えたらFilterValuesを使ってIDをignoreする
func optSyncID(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error {
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

func optRetrieveID(idPtr *uuid.UUID) option {
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
