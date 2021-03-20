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
	id       uuid.UUID `json:"eventId"`
	name     string    `json:"name"`
	duration duration
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

	res := make([]*eventResponse, len(events))
	for _, event := range events {
		res = append(res, &eventResponse{
			id:   event.ID,
			name: event.Name,
			duration: duration{
				since: event.TimeStart,
				until: event.TimeEnd,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

type eventDetailResponse struct {
	eventResponse
	description string `json:"description"`
	place       string `json:"place"`
	hostName    []*userResponse
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
	userRes := make([]*userResponse, len(event.HostName))
	for _, user := range event.HostName {
		userRes = append(userRes, &userResponse{
			id:       user.Id,
			name:     user.Name,
			realName: user.RealName,
		},
		)
	}

	res := &eventDetailResponse{
		eventResponse: eventResponse{
			id:   event.ID,
			name: event.Name,
			duration: duration{
				since: event.TimeStart,
				until: event.TimeEnd,
			},
		},
		description: event.Description,
		place:       event.Place,
		hostName:    userRes,
	}
	return res
}
