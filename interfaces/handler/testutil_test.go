package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

func SetupTestHandlers(t *testing.T, ctrl *gomock.Controller) handler.TestHandlers {
	t.Helper()
	testHandlers := handler.SetupTestApi(ctrl)

	return testHandlers
}

func doRequest(t *testing.T, api handler.API, method, path string, reqBody interface{}, resBody interface{}) (int, *httptest.ResponseRecorder) {
	t.Helper()

	req := httptest.NewRequest(method, path, requestEncode(t, reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()

	infrastructure.Setup(e, api)
	e.ServeHTTP(rec, req)

	// ここ決め打ちじゃないほうが良いかもしれない
	if (rec.Code == http.StatusOK || rec.Code == http.StatusCreated) && !(resBody == nil || reflect.ValueOf(resBody).IsNil()) {
		responseDecode(t, rec, resBody)
	}

	return rec.Code, rec
}

func requestEncode(t *testing.T, body interface{}) *strings.Reader {
	t.Helper()

	b, err := json.Marshal(body)
	require.NoError(t, err)

	return strings.NewReader(string(b))
}

func responseDecode(t *testing.T, rec *httptest.ResponseRecorder, i interface{}) {
	t.Helper()

	err := json.Unmarshal(rec.Body.Bytes(), i)
	require.NoError(t, err)
}

// FIXME: 暫定対処
func ptr(t *testing.T, s string) *string {
	t.Helper()

	return &s
}
