package handler_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
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

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func doRequest(t *testing.T, api handler.API, method, path string, reqBody interface{}, resBody interface{}) (int, *httptest.ResponseRecorder) {
	t.Helper()

	req := httptest.NewRequest(method, path, requestEncode(t, reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	v := validator.New()

	if err := v.RegisterValidation("is-uuid", handler.IsValidUUID); err != nil {
		log.Fatal(err)
	}
	e.Validator = &Validator{
		validator: v,
	}

	infrastructure.Setup(e, api)
	e.ServeHTTP(rec, req)

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
