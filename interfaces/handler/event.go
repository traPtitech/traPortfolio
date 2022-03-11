package handler

import (
	"net/http"
	"time"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventIDInPath struct {
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
		return ConvertError(err)
	}

	res := make([]Event, len(events))
	for i, v := range events {
		res[i] = newEvent(v.ID, v.Name, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

// GetByID GET /events/:eventID
func (h *EventHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
	req := EventIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return ConvertError(err)
	}

	ctx := c.Request().Context()
	event, err := h.srv.GetEventByID(ctx, req.EventID)
	if err != nil {
		return ConvertError(err)
	}

	hostname := make([]User, len(event.HostName))
	for i, v := range event.HostName {
		hostname[i] = newUser(v.ID, v.Name, v.RealName)
	}

	return c.JSON(http.StatusOK, newEventDetail(
		newEvent(event.ID, event.Name, event.TimeStart, event.TimeEnd),
		event.Description,
		event.Level,
		hostname,
		event.Place,
	))
}

// PatchEvent PATCH /events/:eventID
func (h *EventHandler) PatchEvent(_c echo.Context) error {
	c := Context{_c}
	req := struct {
		EventIDInPath
		EditEventJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return ConvertError(err)
	}

	ctx := c.Request().Context()
	patchReq := repository.UpdateEventLevelArg{
		Level: domain.EventLevel(*req.EventLevel),
	}

	if err := h.srv.UpdateEventLevel(ctx, req.EventID, &patchReq); err != nil {
		return ConvertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func newEvent(id uuid.UUID, name string, since time.Time, until time.Time) Event {
	return Event{
		Id:   id,
		Name: name,
		Duration: Duration{
			Since: since,
			Until: &until,
		},
	}
}

func newEventDetail(event Event, description string, eventLevel domain.EventLevel, hostname []User, place string) EventDetail {
	return EventDetail{
		Event:       event,
		Description: description,
		EventLevel:  EventLevel(eventLevel),
		Hostname:    hostname,
		Place:       place,
	}
}
