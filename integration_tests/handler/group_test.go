package handler

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
)

// GetGroups GET /groups
func TestGetGroups(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statusCode int
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			[]schema.Group{
				{
					Id:   mockdata.GroupID1(),
					Name: mockdata.HMockGroups[0].Name,
				},
			},
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Group.GetGroups), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetGroup GET /groups/:groupID
func TestGetGroup(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		statusCode int
		groupID    uuid.UUID
		want       interface{}
	}{
		"200": {
			statusCode: http.StatusOK,
			groupID:    mockdata.GroupID1(),
			want:       mockdata.HMockGroups[0],
		},
		"400 invalid userID": {
			statusCode: http.StatusBadRequest,
			groupID:    uuid.Nil,
			want:       httpError(t, "Bad Request: nil id"),
		},
		"404": {
			statusCode: http.StatusNotFound,
			groupID:    random.UUID(),
			want:       httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Group.GetGroup, tt.groupID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
