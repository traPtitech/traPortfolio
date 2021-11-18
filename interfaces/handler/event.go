package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventHandler struct {
	srv service.EventService
}

// EventResponse Portfolioのレスポンスで使うイベント情報
type EventResponse struct {
	ID       uuid.UUID `json:"eventId"`
	Name     string    `json:"name"`
	Duration Duration
}

// NewEventHandler creates a EventHandler
func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{service}
}

// GetAll GET /events
func (h *EventHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	events, err := h.srv.GetEvents(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]*EventResponse, 0, len(events))
	for _, event := range events {
		res = append(res, &EventResponse{
			ID:   event.ID,
			Name: event.Name,
			Duration: Duration{
				Since: event.TimeStart,
				Until: event.TimeEnd,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

type eventParam struct {
	EventID uuid.UUID `param:"eventID" validate:"is-uuid"`
}

type EventDetailResponse struct {
	EventResponse
	Description string            `json:"description"`
	Place       string            `json:"place"`
	HostName    []*UserResponse   `json:"hostname"`
	EventLevel  domain.EventLevel `json:"eventLevel"`
}

// GetByID GET /events/:eventID
func (h *EventHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
	req := eventParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	event, err := h.srv.GetEventByID(ctx, req.EventID)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, formatUserDetail(event))
}

type EditEventRequest struct {
	EventID    uuid.UUID `param:"eventID" validate:"is-uuid"`
	EventLevel *domain.EventLevel
}

// PatchEvent PATCH /events/:eventID
func (h *EventHandler) PatchEvent(_c echo.Context) error {
	c := Context{_c}
	req := &EditEventRequest{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	patchReq := repository.UpdateEventLevelArg{
		Level: *req.EventLevel,
	}

	if err := h.srv.UpdateEventLevel(ctx, req.EventID, &patchReq); err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func formatUserDetail(event *domain.EventDetail) *EventDetailResponse {
	userRes := make([]*UserResponse, 0, len(event.HostName))
	for _, user := range event.HostName {
		userRes = append(userRes, &UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			RealName: user.RealName,
		},
		)
	}

	res := &EventDetailResponse{
		EventResponse: EventResponse{
			ID:   event.ID,
			Name: event.Name,
			Duration: Duration{
				Since: event.TimeStart,
				Until: event.TimeEnd,
			},
		},
		Description: event.Description,
		Place:       event.Place,
		HostName:    userRes,
		EventLevel:  event.Level,
	}
	return res
}
