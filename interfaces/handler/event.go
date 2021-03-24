package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/event_service"
)

type EventHandler struct {
	service service.EventService
}

// EventResponse Portfolioのレスポンスで使うイベント情報
type eventResponse struct {
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
	events, err := h.service.GetEvents(ctx)
	if err != nil {
		return err
	}

	res := make([]*eventResponse, 0, len(events))
	for _, event := range events {
		res = append(res, &eventResponse{
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

type eventDetailResponse struct {
	eventResponse
	Description string `json:"description"`
	Place       string `json:"place"`
	HostName    []*userResponse
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
	ctx := c.Request().Context()
	event, err := h.service.GetEventByID(ctx, id)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, formatUserDetail(event))
}

func formatUserDetail(event *domain.EventDetail) *eventDetailResponse {
	userRes := make([]*userResponse, 0, len(event.HostName))
	for _, user := range event.HostName {
		userRes = append(userRes, &userResponse{
			ID:       user.ID,
			Name:     user.Name,
			RealName: user.RealName,
		},
		)
	}

	res := &eventDetailResponse{
		eventResponse: eventResponse{
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
	}
	return res
}
