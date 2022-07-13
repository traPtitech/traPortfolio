package handler

import (
	"net/http"
	"time"

	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/optional"

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

// GetEvents GET /events
func (h *EventHandler) GetEvents(c echo.Context) error {
	ctx := c.Request().Context()
	events, err := h.srv.GetEvents(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Event, len(events))
	for i, v := range events {
		res[i] = newEvent(v.ID, v.Name, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

// GetEvent GET /events/:eventID
func (h *EventHandler) GetEvent(_c echo.Context) error {
	c := _c.(*Context)
	req := EventIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	event, err := h.srv.GetEventByID(ctx, req.EventID)
	if err != nil {
		return convertError(err)
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

// EditEvent PATCH /events/:eventID
func (h *EventHandler) EditEvent(_c echo.Context) error {
	c := _c.(*Context)
	req := struct {
		EventIDInPath
		EditEventJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	patchReq := repository.UpdateEventLevelArgs{
		Level: optional.NewUint((uint)(*req.EventLevel), true),
	}

	if err := h.srv.UpdateEventLevel(ctx, req.EventID, &patchReq); err != nil {
		return convertError(err)
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
		Description: description,
		Duration: Duration{
			Since: event.Duration.Since,
			Until: event.Duration.Until,
		},
		EventLevel: EventLevel(eventLevel),
		Hostname:   hostname,
		Id:         event.Id,
		Name:       event.Name,
		Place:      place,
	}
}
