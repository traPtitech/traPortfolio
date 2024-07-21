package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository/mock_repository"
)

type MockRepository struct {
	user    *mock_repository.MockUserRepository
	event   *mock_repository.MockEventRepository
	contest *mock_repository.MockContestRepository
	group   *mock_repository.MockGroupRepository
	project *mock_repository.MockProjectRepository
}

func doRequest(t *testing.T, api API, method, path string, reqBody interface{}, resBody interface{}) (int, *httptest.ResponseRecorder) {
	t.Helper()

	return doRequestWithHeader(t, api, method, path, reqBody, resBody, nil)
}

// TODO: merge with doRequest
func doRequestWithHeader(t *testing.T, api API, method, path string, reqBody interface{}, resBody interface{}, header map[string]string) (int, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, requestEncode(t, reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()

	e := echo.New()

	err := Setup(false, e, api)
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	return strings.NewReader(string(b))
}

func responseDecode(t *testing.T, rec *httptest.ResponseRecorder, i interface{}) {
	t.Helper()

	err := json.Unmarshal(rec.Body.Bytes(), i)
	assert.NoError(t, err)
}

// FIXME: 暫定対処
func ptr(t *testing.T, s string) *string {
	t.Helper()

	return &s
}

type anyCtx struct{}

func (anyCtx) Matches(v interface{}) bool {
	_, ok := v.(context.Context)
	return ok
}

func (anyCtx) String() string {
	return "is Context"
}
