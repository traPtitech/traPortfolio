package handler

import (
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/contest_service"
)

type ContestHandler struct {
	repo    repository.ContestRepository
	service service.ContestService
}

// NewEventHandler creates a EventHandler
func NewContestHandler(repo repository.ContestRepository, service service.ContestService) *ContestHandler {
	return &ContestHandler{repo, service}
}

type PostContestRequest struct {
	Name        string `json:"name" validate:"required"`
	Link        string `json:"link" validate:"url"`
	Description string `json:"description"`
	Duration    Duration
}

type PostContestResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Duration
}

func (h *ContestHandler) PostContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := &PostContestRequest{}
	// todo validation
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	createReq := repository.CreateContestArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Since:       req.Duration.Since,
		Until:       req.Duration.Until,
	}

	contest, err := h.service.CreateContest(ctx, &createReq)
	if err != nil {
		return err
	}
	res := PostContestResponse{
		ID:   contest.ID,
		Name: contest.Name,
		Duration: Duration{
			Since: contest.Since,
			Until: contest.Until,
		},
	}
	return c.JSON(http.StatusCreated, res)
}

type PatchContestRequest struct {
	Name        optional.String `json:"name"`
	Link        optional.String `json:"link"`
	Description optional.String `json:"description"`
	Duration    OptionalDuration
}

func (h *ContestHandler) PatchContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("contestID")
	id := uuid.FromStringOrNil(_id)
	req := &PatchContestRequest{}
	// todo validation
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	patchReq := repository.UpdateContestArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Since:       req.Duration.Since,
		Until:       req.Duration.Until,
	}

	err = h.service.UpdateContest(ctx, id, &patchReq)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
