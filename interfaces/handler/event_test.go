package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func setupEventMock(t *testing.T) (*mock_service.MockEventService, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockEventService(ctrl)
	api := NewAPI(nil, nil, nil, NewEventHandler(s), nil, nil)

	return s, api
}

func TestEventHandler_GetAll(t *testing.T) {

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockEventService) (hres []*Event, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockEventService) (hres []*Event, path string) {

				casenum := 2
				repoEvents := []*domain.Event{}
				hresEvents := []*Event{}

				for i := 0; i < casenum; i++ {
					since, until := random.SinceAndUntil()
					revent := domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(),
						TimeStart: since,
						TimeEnd:   until,
					}
					hevent := Event{
						Id:   revent.ID,
						Name: revent.Name,
						Duration: Duration{
							Since: revent.TimeStart,
							Until: &revent.TimeEnd,
						},
					}

					repoEvents = append(repoEvents, &revent)
					hresEvents = append(hresEvents, &hevent)

				}

				s.EXPECT().GetEvents(anyCtx{}).Return(repoEvents, nil)
				return hresEvents, "/api/v1/events"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockEventService) (hres []*Event, path string) {
				s.EXPECT().GetEvents(anyCtx{}).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/events"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupEventMock(t)

			hresEvents, path := tt.setup(s)

			var resBody []*Event
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresEvents, resBody)
		})
	}
}

func TestEventHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockEventService, hostnum int) (hres *EventDetail, eventpath string)
		statusCode int
	}{
		{
			name: "success random",
			setup: func(s *mock_service.MockEventService, hostnum int) (hres *EventDetail, eventpath string) {

				rHost := []*domain.User{}
				hHost := []User{}

				for i := 0; i < hostnum; i++ {
					rhost := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						RealName: random.AlphaNumeric(),
					}
					hhost := User{
						Id:       rhost.ID,
						Name:     rhost.Name,
						RealName: rhost.RealName,
					}

					rHost = append(rHost, &rhost)
					hHost = append(hHost, hhost)

				}

				since, until := random.SinceAndUntil()
				revent := domain.EventDetail{
					Event: domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(),
						TimeStart: since,
						TimeEnd:   until,
					},

					Description: random.AlphaNumeric(),
					Place:       random.AlphaNumeric(),
					Level:       domain.EventLevel(rand.Intn(domain.EventLevelLimit)),
					HostName:    rHost,
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}

				hevent := EventDetail{
					Description: revent.Description,
					Duration: Duration{
						Since: revent.Event.TimeStart,
						Until: &revent.Event.TimeEnd,
					},
					EventLevel: EventLevel(revent.Level),
					Hostname:   hHost,
					Id:         revent.Event.ID,
					Name:       revent.Event.Name,
					Place:      revent.Place,
				}

				repoEvent := &revent
				hresEvent := &hevent

				s.EXPECT().GetEventByID(anyCtx{}, revent.Event.ID).Return(repoEvent, nil)
				path := fmt.Sprintf("/api/v1/events/%s", revent.Event.ID)
				return hresEvent, path
			},
			statusCode: http.StatusOK,
		},

		{
			name: "internal error",
			setup: func(s *mock_service.MockEventService, hostnum int) (hres *EventDetail, eventpath string) {
				id := random.UUID()
				s.EXPECT().GetEventByID(anyCtx{}, id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/events/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupEventMock(t)

			casenum := []int{1, 2, 32}
			var resBody *EventDetail

			for _, testcase := range casenum {
				hresEvent, eventpath := tt.setup(s, testcase)

				statusCode, _ := doRequest(t, api, http.MethodGet, eventpath, nil, &resBody)

				// Assertion
				assert.Equal(t, tt.statusCode, statusCode)
				assert.Equal(t, hresEvent, resBody)
			}
		})
	}
}

func TestEventHandler_PatchEvent(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockEventService) (reqBody *EditEventRequest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockEventService) (*EditEventRequest, string) {

				eventID := random.UUID()
				eventLevelUint8 := (uint8)(rand.Intn(domain.EventLevelLimit))
				eventLevelHandler := EventLevel(eventLevelUint8)
				//eventLevelDomain := domain.EventLevel(eventLevelUint)

				reqBody := &EditEventRequest{
					EventLevel: &eventLevelHandler,
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.NewUint8((eventLevelUint8), true),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				s.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Conflict",
			setup: func(s *mock_service.MockEventService) (*EditEventRequest, string) {

				eventID := random.UUID()
				eventLevelUint8 := (uint8)(rand.Intn(domain.EventLevelLimit))
				eventLevelHandler := EventLevel(eventLevelUint8)
				//eventLevelDomain := domain.EventLevel(eventLevelUint)

				reqBody := &EditEventRequest{
					EventLevel: &eventLevelHandler,
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.NewUint8((eventLevelUint8), true),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				s.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrAlreadyExists)
				return reqBody, path
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockEventService) (*EditEventRequest, string) {

				eventID := random.UUID()
				eventLevelUint8 := (uint8)(rand.Intn(domain.EventLevelLimit))
				eventLevelHandler := EventLevel(eventLevelUint8)
				//eventLevelDomain := domain.EventLevel(eventLevelUint)

				reqBody := &EditEventRequest{
					EventLevel: &eventLevelHandler,
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.NewUint8((eventLevelUint8), true),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				s.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrNotFound)
				return reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: bind error",
			setup: func(s *mock_service.MockEventService) (*EditEventRequest, string) {

				eventID := random.UUID()
				eventLevelUint8 := (uint8)(rand.Intn(domain.EventLevelLimit))
				eventLevelHandler := EventLevel(eventLevelUint8)
				//eventLevelDomain := domain.EventLevel(eventLevelUint)

				reqBody := &EditEventRequest{
					EventLevel: &eventLevelHandler,
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.NewUint8((eventLevelUint8), true),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				s.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrBind)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error",
			setup: func(s *mock_service.MockEventService) (*EditEventRequest, string) {

				eventID := random.UUID()
				eventLevelUint8 := (uint8)(rand.Intn(domain.EventLevelLimit))
				eventLevelHandler := EventLevel(eventLevelUint8)
				//eventLevelDomain := domain.EventLevel(eventLevelUint)

				reqBody := &EditEventRequest{
					EventLevel: &eventLevelHandler,
				}

				args := repository.UpdateEventLevelArgs{
					Level: optional.NewUint8((eventLevelUint8), true),
				}

				path := fmt.Sprintf("/api/v1/events/%s", eventID)
				s.EXPECT().UpdateEventLevel(anyCtx{}, eventID, &args).Return(repository.ErrValidate)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupEventMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
