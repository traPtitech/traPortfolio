package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
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
	conf := testutils.GetConfigWithDBName("event_handler_get_events")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvents), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
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
			mockdata.HMockEventDetails[0].Id,
			mockdata.HMockEventDetails[0],
		},
		"400 invalid userID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.ConvertError(t, repository.ErrValidate),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			handler.ConvertError(t, repository.ErrNotFound),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("event_handler_get_event")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvent, tt.eventID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// EditEvent PATCH /events/:eventID
func TestEventHandler_EditEvent(t *testing.T) {
	var eventLevel = handler.EventLevel(rand.Intn(domain.EventLevelLimit))

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		eventID    uuid.UUID
		reqBody    handler.EditEventRequest
		want       interface{} // nil or error
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockEventDetails[0].Id,
			handler.EditEventRequest{
				EventLevel: &eventLevel,
			},
			nil,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("event_handler_edit_event")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var event handler.EventDetail
				res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvent, tt.eventID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &event)) // TODO: ここだけjson.Unmarshalを直接行っているのでスマートではない

				// Update & Assert
				res = testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.Event.EditEvent, tt.eventID), tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)

				// Get updated response & Assert
				if tt.reqBody.EventLevel != nil {
					event.EventLevel = *tt.reqBody.EventLevel
				}
				res = testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Event.GetEvent, tt.eventID), nil)
				testutils.AssertResponse(t, http.StatusOK, event, res)
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
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
                        t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.XXX.GetXXX, tt.userID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
*/
