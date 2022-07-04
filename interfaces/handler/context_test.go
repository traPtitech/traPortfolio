package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContext_BindAndValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		reqBodyStr string
		wantErr    bool
	}{
		"ok: empty":                      {`{}`, false},
		"ok: with name":                  {`{"name": "test"}`, false},
		"ok: with includeSuspended":      {`{"includeSuspended": true}`, false},
		"ok: with irrelevant param":      {`{"name": "test", "a": "a"}`, false},
		"ng(validate): with both params": {`{"name": "test", "includeSuspended": true}`, true},
		"ng(validate): with empty name":  {`{"name": ""}`, true},
		"ng(bind): invalid json":         {`{"name":`, true},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", strings.NewReader(tt.reqBodyStr))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			e.Validator = newValidator(e.Logger)
			ctx := e.NewContext(req, rec)
			c := &Context{ctx}

			err := c.BindAndValidate(&GetUsersParams{})
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
