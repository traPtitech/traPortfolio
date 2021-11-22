package handler_test

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestNewEventHandler(t *testing.T) {
	/*tests := []struct {
		name    string
		setup   func() *EventHandler
		service service.EventService
		want    *EventHandler
	}{
		// TODO: Add test cases.
		{
			name: "Success",
			setup: func() *EventHandler {
				return NewEventHandler(service)
			},
			service: service.EventService{
				event: repository.EventRepository{

				},
				user:  repository.UserRepository{ae},
			},
			want: &EventHandler{service},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//got := NewEventHandler(tt.setup(t.context))
			got := tt.setup()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventHandler() = %v, want %v", got, tt.want)
			}
		})
	}*/
}

func TestEventHandler_GetAll(t *testing.T) {

	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (hres []*handler.EventResponse, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(th *handler.TestHandlers) (hres []*handler.EventResponse, path string) {

				casenum := 2
				repoEvents := []*domain.Event{}
				hresEvents := []*handler.EventResponse{}

				for i := 0; i < casenum; i++ {
					revent := domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					}
					hevent := handler.EventResponse{
						ID:   revent.ID,
						Name: revent.Name,
						Duration: handler.Duration{
							Since: revent.TimeStart,
							Until: revent.TimeEnd,
						},
					}

					repoEvents = append(repoEvents, &revent)
					hresEvents = append(hresEvents, &hevent)

				}

				th.Service.MockEventService.EXPECT().GetEvents(gomock.Any()).Return(repoEvents, nil)
				return hresEvents, "/api/v1/events"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(th *handler.TestHandlers) (hres []*handler.EventResponse, path string) {
				th.Service.MockEventService.EXPECT().GetEvents(gomock.Any()).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/events"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			hresEvents, path := tt.setup(&handlers)

			var resBody []*handler.EventResponse
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresEvents, resBody)
		})
	}
}

func TestEventHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (hres *handler.EventDetailResponse, eventpath string)
		statusCode int
	}{
		{
			name: "success 1",
			setup: func(th *handler.TestHandlers) (hres *handler.EventDetailResponse, eventpath string) {

				//1人
				hostnum := 1
				rHost := []*domain.User{}
				hHost := []*handler.UserResponse{}

				for i := 0; i < hostnum; i++ {
					rhost := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					}
					hhost := handler.UserResponse{
						ID:       rhost.ID,
						Name:     rhost.Name,
						RealName: rhost.RealName,
					}

					rHost = append(rHost, &rhost)
					hHost = append(hHost, &hhost)

				}

				revent := domain.EventDetail{

					Event: domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					},

					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Place:       random.AlphaNumeric(rand.Intn(30) + 1),
					Level:       domain.EventLevel(rand.Intn(domain.EventLevelLimit)),
					HostName:    rHost,
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}

				hevent := handler.EventDetailResponse{
					EventResponse: handler.EventResponse{
						ID:   revent.Event.ID,
						Name: revent.Event.Name,
						Duration: handler.Duration{
							Since: revent.Event.TimeStart,
							Until: revent.Event.TimeEnd,
						},
					},

					Description: revent.Description,
					Place:       revent.Place,
					HostName:    hHost,
					EventLevel:  revent.Level,
				}

				repoEvent := &revent
				hresEvent := &hevent

				th.Service.MockEventService.EXPECT().GetEventByID(gomock.Any(), revent.Event.ID).Return(repoEvent, nil)
				path := fmt.Sprintf("/api/v1/events/%s", revent.Event.ID)
				return hresEvent, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(th *handler.TestHandlers) (hres *handler.EventDetailResponse, eventpath string) {

				//2人
				hostnum := 2
				rHost := []*domain.User{}
				hHost := []*handler.UserResponse{}

				for i := 0; i < hostnum; i++ {
					rhost := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					}
					hhost := handler.UserResponse{
						ID:       rhost.ID,
						Name:     rhost.Name,
						RealName: rhost.RealName,
					}

					rHost = append(rHost, &rhost)
					hHost = append(hHost, &hhost)

				}

				revent := domain.EventDetail{

					Event: domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					},

					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Place:       random.AlphaNumeric(rand.Intn(30) + 1),
					Level:       domain.EventLevel(rand.Intn(domain.EventLevelLimit)),
					HostName:    rHost,
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}

				hevent := handler.EventDetailResponse{
					EventResponse: handler.EventResponse{
						ID:   revent.Event.ID,
						Name: revent.Event.Name,
						Duration: handler.Duration{
							Since: revent.Event.TimeStart,
							Until: revent.Event.TimeEnd,
						},
					},

					Description: revent.Description,
					Place:       revent.Place,
					HostName:    hHost,
					EventLevel:  revent.Level,
				}

				repoEvent := &revent
				hresEvent := &hevent

				th.Service.MockEventService.EXPECT().GetEventByID(gomock.Any(), revent.Event.ID).Return(repoEvent, nil)
				path := fmt.Sprintf("/api/v1/events/%s", revent.Event.ID)
				return hresEvent, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 3",
			setup: func(th *handler.TestHandlers) (hres *handler.EventDetailResponse, eventpath string) {

				//32人
				hostnum := 32
				rHost := []*domain.User{}
				hHost := []*handler.UserResponse{}

				for i := 0; i < hostnum; i++ {
					rhost := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					}
					hhost := handler.UserResponse{
						ID:       rhost.ID,
						Name:     rhost.Name,
						RealName: rhost.RealName,
					}

					rHost = append(rHost, &rhost)
					hHost = append(hHost, &hhost)

				}

				revent := domain.EventDetail{

					Event: domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					},

					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Place:       random.AlphaNumeric(rand.Intn(30) + 1),
					Level:       domain.EventLevel(rand.Intn(domain.EventLevelLimit)),
					HostName:    rHost,
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}

				hevent := handler.EventDetailResponse{
					EventResponse: handler.EventResponse{
						ID:   revent.Event.ID,
						Name: revent.Event.Name,
						Duration: handler.Duration{
							Since: revent.Event.TimeStart,
							Until: revent.Event.TimeEnd,
						},
					},

					Description: revent.Description,
					Place:       revent.Place,
					HostName:    hHost,
					EventLevel:  revent.Level,
				}

				repoEvent := &revent
				hresEvent := &hevent

				th.Service.MockEventService.EXPECT().GetEventByID(gomock.Any(), revent.Event.ID).Return(repoEvent, nil)
				path := fmt.Sprintf("/api/v1/events/%s", revent.Event.ID)
				return hresEvent, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(th *handler.TestHandlers) (hres *handler.EventDetailResponse, eventpath string) {
				id := random.UUID()
				th.Service.MockEventService.EXPECT().GetEventByID(gomock.Any(), id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/events/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			hresEvent, eventpath := tt.setup(&handlers)

			var resBody *handler.EventDetailResponse
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, eventpath, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresEvent, resBody)
		})
	}
}

/*
func TestEventHandler_PatchEvent(t *testing.T) {
	type fields struct {
		srv service.EventService
	}
	type args struct {
		_c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &EventHandler{
				srv: tt.fields.srv,
			}
			if err := h.PatchEvent(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("EventHandler.PatchEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_formatUserDetail(t *testing.T) {
	type args struct {
		event *domain.EventDetail
	}
	tests := []struct {
		name string
		args args
		want *eventDetailResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatUserDetail(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatUserDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
