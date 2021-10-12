package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func SetupTestHandlers(t *testing.T, ctrl *gomock.Controller) TestHandlers {
	t.Helper()
	testHandlers := SetupTestApi(ctrl)

	return testHandlers
}

func doRequest(t *testing.T, handler func(echo.Context) error, method, path string, reqBody interface{}, resBody interface{}) (int, *httptest.ResponseRecorder) {
	t.Helper()

	req := httptest.NewRequest(method, path, requestEncode(t, reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := echo.New().NewContext(req, rec)
	handler(ctx)
	if resBody != nil {
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
