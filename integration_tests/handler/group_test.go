package handler

import (
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
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
			[]handler.Group{
				{
					Id:   mockdata.GroupID(),
					Name: mockdata.HMockGroups[0].Name,
				},
			},
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("group_handler_get_groups")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Group.GetGroups), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
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
			groupID:    mockdata.GroupID(),
			want:       mockdata.HMockGroups[0],
		},
		"400 invalid userID": {
			statusCode: http.StatusBadRequest,
			groupID:    uuid.Nil,
			want:       testutils.HTTPError("bad request: nil id"),
		},
		"404": {
			statusCode: http.StatusNotFound,
			groupID:    random.UUID(),
			want:       testutils.HTTPError("not found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("group_handler_get_group")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Group.GetGroup, tt.groupID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
