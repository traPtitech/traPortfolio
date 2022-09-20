package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

func SetupRoutes(t *testing.T, e *echo.Echo, conf *config.Config) (*handler.API, error) {
	t.Helper()

	db := SetupDB(t, conf.SQLConf())
	if err := mockdata.InsertSampleDataToDB(db); err != nil {
		return nil, err
	}

	api, err := infrastructure.InjectAPIServer(conf)
	if err != nil {
		return nil, err
	}

	if err := handler.Setup(e, api); err != nil {
		return nil, err
	}

	return &api, nil
}

func DoRequest(t *testing.T, e *echo.Echo, method string, path string, body interface{}) *httptest.ResponseRecorder {
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
func DummyUUID() uuid.UUID {
	return random.UUID()
}

func HTTPError(message string) echo.HTTPError {
	return echo.HTTPError{
		Message: message,
	}
}
