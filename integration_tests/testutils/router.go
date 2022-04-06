//go:build integration && db

package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
)

type initDBFunc func(database.SQLHandler) error

func SetupRoutes(t *testing.T, e *echo.Echo, conf *config.Config, f initDBFunc) (*handler.API, error) {
	t.Helper()

	s := conf.SQLConf()
	tr := conf.TraqConf()
	p := conf.PortalConf()
	k := conf.KnoqConf()

	db := SetupDB(t, &s)
	if f != nil {
		if err := f(db); err != nil {
			return nil, err
		}
	}

	api, err := infrastructure.InjectAPIServer(&s, &tr, &p, &k)
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
