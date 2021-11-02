package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type eventIDInPath struct {
	EventID uuid.UUID `param:"eventID" validate:"is-uuid"`
}

type EventHandler struct {
	srv service.EventService
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

	res := make([]*Event, 0, len(events))
	for _, event := range events {
		res = append(res, &Event{
			Id:   event.ID,
			Name: event.Name,
			Duration: Duration{
				Since: event.TimeStart,
				Until: &event.TimeEnd,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

// GetByID GET /events/:eventID
func (h *EventHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
	req := eventIDInPath{}
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

// PatchEvent PATCH /events/:eventID
func (h *EventHandler) PatchEvent(_c echo.Context) error {
	c := Context{_c}
	req := struct {
		eventIDInPath
		EditEventJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	patchReq := repository.UpdateEventArg{
		// Level: *req.EventLevel, // TODO
	}

	if err := h.srv.UpdateEvent(ctx, req.EventID, &patchReq); err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func formatUserDetail(event *domain.EventDetail) *EventDetail {
	userRes := make([]*User, 0, len(event.HostName))
	for _, user := range event.HostName {
		userRes = append(userRes, &User{
			Id:       user.ID,
			Name:     user.Name,
			RealName: &user.RealName,
		},
		)
	}

	res := &EventDetail{
		Event: Event{
			Id:   event.ID,
			Name: event.Name,
			Duration: Duration{
				Since: event.TimeStart,
				Until: &event.TimeEnd,
			},
		},
		Description: event.Description,
		Place:       event.Place,
		// Hostname:    userRes, // TODO
		// EventLevel:  EventLevel(event.Level), // TODO
	}
	return res
}
