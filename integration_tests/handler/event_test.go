package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

// GetEvents GET /events
func TestEventHandler_GetEvents(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockEvents,
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvents), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetEvent GET /events/:eventID
func TestEventHandler_GetEvent(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		eventID    uuid.UUID
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.KnoqEventID1(),
			mockdata.HMockEventDetails[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvent, tt.eventID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// EditEvent PATCH /events/:eventID
func TestEventHandler_EditEvent(t *testing.T) {
	var (
		eventLevel = schema.EventLevel(domain.EventLevelPublic)
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		eventID    uuid.UUID
		reqBody    schema.EditEventRequest
		want       interface{} // nil or error
	}{
		"204": {
			http.StatusNoContent,
			mockdata.KnoqEventID1(),
			schema.EditEventRequest{
				Level: &eventLevel,
			},
			nil,
		},
		"204 without change": {
			http.StatusNoContent,
			mockdata.KnoqEventID3(),
			schema.EditEventRequest{},
			nil,
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var event schema.EventDetail
				res := doRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvent, tt.eventID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &event)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない

				// Update & Assert
				res = doRequest(t, e, http.MethodPatch, e.URL(api.Event.EditEvent, tt.eventID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)

				// Get updated response & Assert
				if tt.reqBody.Level != nil {
					event.Level = *tt.reqBody.Level
				}
				res = doRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvent, tt.eventID), nil)
				assertResponse(t, http.StatusOK, event, res)
			}
		})
	}
}

/*
// GetXXX GET /XXX
func TestGetXXX(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       interface{}
	}{
		// TODO: Add cases
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("xxx_handler_get_xxx")
	api, err := setupRoutes(t, e,)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
                        t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.XXX.GetXXX, tt.userID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
*/
