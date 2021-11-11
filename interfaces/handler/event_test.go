package handler_test

import (
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

				repoEvents := []*domain.Event{
					{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					},
					{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					},
				}

				hresEvents := []*handler.EventResponse{
					{
						ID:   repoEvents[0].ID,
						Name: repoEvents[0].Name,
						Duration: handler.Duration{
							Since: repoEvents[0].TimeStart,
							Until: repoEvents[0].TimeEnd,
						},
					},
					{
						ID:   repoEvents[1].ID,
						Name: repoEvents[1].Name,
						Duration: handler.Duration{
							Since: repoEvents[1].TimeStart,
							Until: repoEvents[1].TimeEnd,
						},
					},
				}

				th.Service.MockEventService.EXPECT().GetEvents(gomock.Any()).Return(repoEvents, nil)
				return hresEvents, "/api/v1/events"
			},
			statusCode: http.StatusOK,
		},
		/*{
			name: "internal error",
			setup: func(th *handler.TestHandlers, want []*domain.Event) string {
				th.Repository.MockEventRepository.EXPECT().GetEvents().Return(nil, err?)
				return "/api/v1/events"
			},
			statusCode: 500,
		},*/
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

/*func TestEventHandler_GetByID(t *testing.T) {
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
			if err := h.GetByID(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("EventHandler.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
