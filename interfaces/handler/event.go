package handler

import (
	"net/http"
	"time"

	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

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
		return err
	}

	res := make([]schema.Event, len(events))
	for i, v := range events {
		res[i] = newEvent(v.ID, v.Name, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

// GetEvent GET /events/:eventID
func (h *EventHandler) GetEvent(c echo.Context) error {
	eventID, err := getID(c, keyEventID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	event, err := h.srv.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	hostname := make([]schema.User, len(event.HostName))
	for i, v := range event.HostName {
		hostname[i] = newUser(v.ID, v.Name, v.RealName())
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
func (h *EventHandler) EditEvent(c echo.Context) error {
	eventID, err := getID(c, keyEventID)
	if err != nil {
		return err
	}

	req := schema.EditEventJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	patchReq := repository.UpdateEventLevelArgs{
		Level: optional.FromPtr((*domain.EventLevel)(req.EventLevel)),
	}

	if err := h.srv.UpdateEventLevel(ctx, eventID, &patchReq); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func newEvent(id uuid.UUID, name string, since time.Time, until time.Time) schema.Event {
	return schema.Event{
		Id:   id,
		Name: name,
		Duration: schema.Duration{
			Since: since,
			Until: &until,
		},
	}
}

func newEventDetail(event schema.Event, description string, eventLevel domain.EventLevel, hostname []schema.User, place string) schema.EventDetail {
	return schema.EventDetail{
		Description: description,
		Duration: schema.Duration{
			Since: event.Duration.Since,
			Until: event.Duration.Until,
		},
		EventLevel: schema.EventLevel(eventLevel),
		Hostname:   hostname,
		Id:         event.Id,
		Name:       event.Name,
		Place:      place,
	}
}
