package handler

import (
	"net/http"
	"time"

	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"

	"github.com/traPtitech/traPortfolio/internal/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
)

type EventHandler struct {
	event repository.EventRepository
	user  repository.UserRepository
}

// NewEventHandler creates a EventHandler
func NewEventHandler(event repository.EventRepository, user repository.UserRepository) *EventHandler {
	return &EventHandler{event, user}
}

// GetEvents GET /events
func (h *EventHandler) GetEvents(c echo.Context) error {
	ctx := c.Request().Context()
	events, err := h.event.GetEvents(ctx)
	if err != nil {
		return err
	}

	res := make([]schema.Event, len(events))
	for i, v := range events {
		res[i] = newEvent(v.ID, v.Name, v.Level, v.TimeStart, v.TimeEnd)
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
	event, err := h.event.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}
	{
		// hostnameの詳細を取得
		users, err := h.user.GetUsers(ctx, &repository.GetUsersArgs{}) // TODO: IncludeSuspendedをtrueにするか考える
		if err != nil {
			return err
		}

		umap := make(map[uuid.UUID]*domain.User)
		for _, u := range users {
			umap[u.ID] = u
		}

		for i, v := range event.HostName {
			if u, ok := umap[v.ID]; ok {
				event.HostName[i] = u
			}
		}
	}
	hostname := make([]schema.User, len(event.HostName))
	for i, v := range event.HostName {
		hostname[i] = newUser(v.ID, v.Name, v.RealName())
	}

	return c.JSON(http.StatusOK, newEventDetail(
		newEvent(event.ID, event.Name, event.Level, event.TimeStart, event.TimeEnd),
		event.Description,
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

	req := schema.EditEventRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	patchReq := repository.UpdateEventLevelArgs{
		Level: optional.FromPtr((*domain.EventLevel)(req.Level)),
	}

	if err := h.event.UpdateEventLevel(ctx, eventID, &patchReq); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func newEvent(id uuid.UUID, name string, level domain.EventLevel, since time.Time, until time.Time) schema.Event {
	return schema.Event{
		Id:    id,
		Name:  name,
		Level: schema.EventLevel(level),
		Duration: schema.Duration{
			Since: since,
			Until: &until,
		},
	}
}

func newEventDetail(event schema.Event, description string, hostname []schema.User, place string) schema.EventDetail {
	return schema.EventDetail{
		Description: description,
		Duration: schema.Duration{
			Since: event.Duration.Since,
			Until: event.Duration.Until,
		},
		Level:    event.Level,
		Hostname: hostname,
		Id:       event.Id,
		Name:     event.Name,
		Place:    place,
	}
}
