package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ContestHandler struct {
	srv service.ContestService
}

// NewEventHandler creates a EventHandler
func NewContestHandler(service service.ContestService) *ContestHandler {
	return &ContestHandler{service}
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

	contest, err := h.srv.CreateContest(ctx, &createReq)
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

	err = h.srv.UpdateContest(ctx, id, &patchReq)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type PostContestTeamRequest struct {
	Name        string `json:"name"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

type PostContestTeamResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Result string    `json:"result"`
}

func (h *ContestHandler) PostContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("contestID")
	id := uuid.FromStringOrNil(_id)
	req := &PostContestTeamRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.CreateContestTeamArgs{
		Name:        req.Name,
		Result:      req.Result,
		Link:        req.Link,
		Description: req.Description,
	}
	contestTeam, err := h.srv.CreateContestTeam(ctx, id, &args)
	if err != nil {
		return err
	}

	res := &PostContestTeamRequest{
		Name:        contestTeam.Name,
		Link:        contestTeam.Link,
		Description: contestTeam.Description,
		Result:      contestTeam.Result,
	}
	return c.JSON(http.StatusCreated, res)
}

type PatchContestTeamRequest struct {
	Name        optional.String `json:"name"`
	Link        optional.String `json:"link" validate:"url"`
	Description optional.String `json:"description"`
	Result      optional.String `json:"result"`
}

func (h *ContestHandler) PatchContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	_ = uuid.FromStringOrNil(c.Param("contestID"))
	teamID := uuid.FromStringOrNil(c.Param("teamID"))
	req := &PatchContestTeamRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.UpdateContestTeamArgs{
		Name:        req.Name,
		Result:      req.Result,
		Link:        req.Link,
		Description: req.Description,
	}

	err = h.srv.UpdateContestTeam(ctx, teamID, &args)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type PutContestTeamMember struct {
	Members []uuid.UUID `json:"members"`
}

func (h *ContestHandler) PutContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	_ = uuid.FromStringOrNil(c.Param("contestID"))
	teamID := uuid.FromStringOrNil(c.Param("teamID"))
	req := &PutContestTeamMember{}
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.srv.AddContestTeamMember(ctx, teamID, req.Members)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ContestHandler) DeleteContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	_ = uuid.FromStringOrNil(c.Param("contestID"))
	teamID := uuid.FromStringOrNil(c.Param("teamID"))
	req := &PutContestTeamMember{}
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.srv.DeleteContestTeamMember(ctx, teamID, req.Members)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
