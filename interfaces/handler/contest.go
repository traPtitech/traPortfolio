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

// NewContestHandler creates a ContestHandler
func NewContestHandler(service service.ContestService) *ContestHandler {
	return &ContestHandler{service}
}

type PostContestRequest struct {
	Name        string `json:"name" validate:"required"`
	Link        string `json:"link" validate:"url"`
	Description string `json:"description"`
	Duration    Duration
}

type ContestResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Duration Duration  `json:"duration"`
}

type ContestDetailResponse struct {
	ContestResponse
	Link        string                 `json:"link"`
	Description string                 `json:"description"`
	Teams       []*ContestTeamResponse `json:"teams"`
}

type ContestTeamResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Result string    `json:"result"`
}

type ContestTeamDetailResponse struct {
	ContestTeamResponse
	Link        string          `json:"link"`
	Description string          `json:"description"`
	Members     []*userResponse `json:"members"`
}

// GetContests GET /contests
func (h *ContestHandler) GetContests(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	contests, err := h.srv.GetContests(ctx)
	if err != nil {
		return convertError(err)
	}
	res := make([]*ContestResponse, 0, len(contests))
	for _, v := range contests {
		res = append(res, &ContestResponse{
			ID:   v.ID,
			Name: v.Name,
			Duration: Duration{
				Since: v.TimeStart,
				Until: v.TimeEnd,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *ContestHandler) GetContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("contestID")
	id := uuid.FromStringOrNil(_id)

	contest, err := h.srv.GetContest(ctx, id)

	teams := make([]*ContestTeamResponse, 0, len(contest.Teams))
	for _, v := range contest.Teams {
		teams = append(teams, &ContestTeamResponse{
			ID:     v.ID,
			Name:   v.Name,
			Result: v.Result,
		})
	}

	res := &ContestDetailResponse{
		ContestResponse: ContestResponse{
			ID:   contest.ID,
			Name: contest.Name,
			Duration: Duration{
				Since: contest.TimeStart,
				Until: contest.TimeEnd,
			},
		},
		Link:        contest.Link,
		Description: contest.Description,
		Teams:       teams,
	}

	if err != nil {
		return convertError(err)
	}
	return c.JSON(http.StatusOK, res)
}

// PostContest POST /contests
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
	res := ContestResponse{
		ID:   contest.ID,
		Name: contest.Name,
		Duration: Duration{
			Since: contest.TimeStart,
			Until: contest.TimeEnd,
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

// PatchContest PATCH /contests/:contestID
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

// DeleteContest DELETE /contests/:contestID
func (h *ContestHandler) DeleteContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("contestID")
	id := uuid.FromStringOrNil(_id)
	err := h.srv.DeleteContest(ctx, id)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
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

// GetContestTeams GET /contests/:contestID/teams
func (h *ContestHandler) GetContestTeams(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("contestID")
	id := uuid.FromStringOrNil(_id)
	contestTeams, err := h.srv.GetContestTeams(ctx, id)
	if err != nil {
		return err
	}

	res := make([]*ContestTeamResponse, 0, len(contestTeams))
	for _, v := range contestTeams {
		ct := &ContestTeamResponse{
			ID:     v.ID,
			Name:   v.Name,
			Result: v.Result,
		}
		res = append(res, ct)
	}
	return c.JSON(http.StatusOK, res)
}

// GetContestTeams GET /contests/:contestID/teams/:teamID
func (h *ContestHandler) GetContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("contestID")
	contestID := uuid.FromStringOrNil(_id)
	_id = c.Param("teamID")
	teamID := uuid.FromStringOrNil(_id)
	contestTeam, err := h.srv.GetContestTeam(ctx, contestID, teamID)
	if err != nil {
		return err
	}

	members := make([]*userResponse, 0, len(contestTeam.Members))
	for _, user := range contestTeam.Members {
		members = append(members, &userResponse{
			ID:       user.ID,
			Name:     user.Name,
			RealName: user.RealName,
		})	
	}

	res := &ContestTeamDetailResponse{
		ContestTeamResponse: ContestTeamResponse{
			ID:     contestTeam.ID,
			Name:   contestTeam.Name,
			Result: contestTeam.Result,
		},
		Link:        contestTeam.Link,
		Description: contestTeam.Description,
		Members:     members,
	}
	return c.JSON(http.StatusOK, res)
}

// PostContestTeam POST /contests/:contestID/teams
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

// PatchContestTeam PATCH /contests/:contestID
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

// GetContestTeamMember GET /contests/{contestId}/teams/{teamId}/members
func (h *ContestHandler) GetContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	contestID := uuid.FromStringOrNil(c.Param("contestID"))
	teamID := uuid.FromStringOrNil(c.Param("teamID"))

	users, err := h.srv.GetContestTeamMember(ctx, contestID, teamID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*userResponse, 0, len(users))
	for _, v := range users {
		res = append(res, &userResponse{
			ID:       v.ID,
			Name:     v.Name,
			RealName: v.RealName,
		})
	}
	return c.JSON(http.StatusOK, res)
}

type PostContestTeamMember struct {
	Members []uuid.UUID `json:"members" validate:"required"`
}

// PostContestTeamMember POST /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) PostContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	_ = uuid.FromStringOrNil(c.Param("contestID"))
	teamID := uuid.FromStringOrNil(c.Param("teamID"))
	req := &PostContestTeamMember{}
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

// DeleteContestTeamMember DELETE /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) DeleteContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	_ = uuid.FromStringOrNil(c.Param("contestID"))
	teamID := uuid.FromStringOrNil(c.Param("teamID"))
	req := &PostContestTeamMember{}
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
