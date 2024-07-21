package handler

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository/mock_repository"
)

func setupEventMock(t *testing.T) (MockRepository, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	event := mock_repository.NewMockEventRepository(ctrl)
	user := mock_repository.NewMockUserRepository(ctrl)
	mr := MockRepository{user: user, event: event}
	api := NewAPI(nil, nil, nil, NewEventHandler(event, user), nil, nil)

	return mr, api
}

func TestEventHandler_GetEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres []*schema.Event, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				casenum := 2
				repoEvents := []*domain.Event{}
				hresEvents := []*schema.Event{}

				for range casenum {
					since, until := random.SinceAndUntil()
					revent := domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(),
						TimeStart: since,
						TimeEnd:   until,
					}
					hevent := schema.Event{
						Id:   revent.ID,
						Name: revent.Name,
						Duration: schema.Duration{
							Since: revent.TimeStart,
							Until: &revent.TimeEnd,
						},
					}

					repoEvents = append(repoEvents, &revent)
					hresEvents = append(hresEvents, &hevent)
				}

				mr.event.EXPECT().GetEvents(anyCtx{}).Return(repoEvents, nil)
				return hresEvents, "/api/v1/events"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				mr.event.EXPECT().GetEvents(anyCtx{}).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/events"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupEventMock(t)

			hresEvents, path := tt.setup(mr)

			var resBody []*schema.Event
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresEvents, resBody)
		})
	}
}

func TestEventHandler_GetEvent(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(mr MockRepository, hostnum int) (hres *schema.EventDetail, eventpath string)
		statusCode int
	}{
		{
			name: "success random",
			setup: func(mr MockRepository, hostnum int) (hres *schema.EventDetail, eventpath string) {
				rHost := []*domain.User{}
				hHost := []schema.User{}

				for range hostnum {
					rhost := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())
					hhost := schema.User{
						Id:       rhost.ID,
						Name:     rhost.Name,
						RealName: rhost.RealName(),
					}

					rHost = append(rHost, rhost)
					hHost = append(hHost, hhost)
				}

				since, until := random.SinceAndUntil()
				revent := domain.EventDetail{
					Event: domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(),
						Level:     rand.N(domain.EventLevelLimit),
						TimeStart: since,
						TimeEnd:   until,
					},
					Description: random.AlphaNumeric(),
					Place:       random.AlphaNumeric(),
					HostName:    rHost,
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}

				hevent := schema.EventDetail{
					Description: revent.Description,
					Duration: schema.Duration{
						Since: revent.Event.TimeStart,
						Until: &revent.Event.TimeEnd,
					},
					Level:    schema.EventLevel(revent.Level),
					Hostname: hHost,
					Id:       revent.Event.ID,
					Name:     revent.Event.Name,
					Place:    revent.Place,
				}

				repoEvent := &revent
				hresEvent := &hevent

				mr.event.EXPECT().GetEvent(anyCtx{}, revent.Event.ID).Return(repoEvent, nil)
				mr.user.EXPECT().GetUsers(anyCtx{}, &repository.GetUsersArgs{}).Return(rHost, nil)
				path := fmt.Sprintf("/api/v1/events/%s", revent.Event.ID)
				return hresEvent, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "BadRequest: Invalid event ID",
			setup: func(_ MockRepository, hostnum int) (hres *schema.EventDetail, eventpath string) {
				return nil, fmt.Sprintf("/api/v1/events/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository, _ int) (hres *schema.EventDetail, eventpath string) {
				id := random.UUID()
				mr.event.EXPECT().GetEvent(anyCtx{}, id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/events/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mr, api := setupEventMock(t)

			casenum := []int{1, 2, 32}
			var resBody *schema.EventDetail

			for _, testcase := range casenum {
				hresEvent, eventpath := tt.setup(mr, testcase)

				statusCode, _ := doRequest(t, api, http.MethodGet, eventpath, nil, &resBody)

				// Assertion
				assert.Equal(t, tt.statusCode, statusCode)
				assert.Equal(t, hresEvent, resBody)
			}
		})
	}
}

func TestEventHandler_EditEvent(t *testing.T) {
	hLevel := func(l domain.EventLevel) *schema.EventLevel {
		r := schema.EventLevel(l)
		return &r
	}

	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (reqBody *schema.EditEventRequest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (*schema.EditEventRequest, string) {
				eventID := random.UUID()
				eventLevel := rand.N(domain.EventLevelLimit)

				reqBody := &schema.EditEventRequest{
					Level: hLevel(eventLevel),
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.From(eventLevel),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				mr.event.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid event ID",
			setup: func(_ MockRepository) (*schema.EditEventRequest, string) {
				return nil, fmt.Sprintf("/api/v1/events/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Conflict",
			setup: func(mr MockRepository) (*schema.EditEventRequest, string) {
				eventID := random.UUID()
				eventLevel := rand.N(domain.EventLevelLimit)

				reqBody := &schema.EditEventRequest{
					Level: hLevel(eventLevel),
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.From(eventLevel),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				mr.event.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrAlreadyExists)
				return reqBody, path
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "Not Found",
			setup: func(mr MockRepository) (*schema.EditEventRequest, string) {
				eventID := random.UUID()
				eventLevel := rand.N(domain.EventLevelLimit)

				reqBody := &schema.EditEventRequest{
					Level: hLevel(eventLevel),
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.From(eventLevel),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				mr.event.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrNotFound)
				return reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: bind error",
			setup: func(mr MockRepository) (*schema.EditEventRequest, string) {
				eventID := random.UUID()
				eventLevel := rand.N(domain.EventLevelLimit)

				reqBody := &schema.EditEventRequest{
					Level: hLevel(eventLevel),
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.From(eventLevel),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				mr.event.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrBind)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: too large level",
			setup: func(_ MockRepository) (*schema.EditEventRequest, string) {
				eventID := random.UUID()
				eventLevel := schema.EventLevel(domain.EventLevelLimit)

				reqBody := &schema.EditEventRequest{
					Level: &eventLevel,
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)

				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupEventMock(t)

			reqBody, path := tt.setup(mr)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
