package handler

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/event_service"
)

type EventHandler struct {
	repo    repository.EventRepository
	service service.EventService
}

// NewEventHandler creates a EventHandler
func NewEventHandler(repo repository.EventRepository, service service.EventService) *EventHandler {
	return &EventHandler{repo, service}
}

// GetAll GET /events
func (h *EventHandler) GetAll(c echo.Context) error {
	ctx := context.Background()
	events, err := h.service.GetEvents(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, events)
}

// GetByID GET /events/:eventID
func (h *EventHandler) GetByID(c echo.Context) error {
	_id := c.Param("eventID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "event id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}
	ctx := context.Background()
	event, err := h.service.GetEventByID(ctx, id)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, event)
}
