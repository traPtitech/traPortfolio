package handler

import (
	"net/http"

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
	req := EventIDInPath{}
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
		EventIDInPath
		EditEventJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	patchReq := repository.UpdateEventLevelArg{
		Level: domain.EventLevel(*req.EventLevel),
	}

	if err := h.srv.UpdateEventLevel(ctx, req.EventID, &patchReq); err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func formatUserDetail(event *domain.EventDetail) *EventDetail {
	userRes := make([]User, len(event.HostName))
	for i, user := range event.HostName {
		userRes[i] = User{
			Id:       user.ID,
			Name:     user.Name,
			RealName: &user.RealName,
		}
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
		Hostname:    userRes,
		EventLevel:  EventLevel(event.Level),
	}
	return res
}
