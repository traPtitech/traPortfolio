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

type contestParam struct {
	ContestID uuid.UUID `param:"contestID" validate:"is-uuid"`
}

type teamParams struct {
	ContestID uuid.UUID `param:"contestID" validate:"is-uuid"`
	TeamID    uuid.UUID `param:"teamID" validate:"is-uuid"`
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
	req := contestParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	contest, err := h.srv.GetContest(ctx, req.ContestID)
	if err != nil {
		return convertError(err)
	}

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

	return c.JSON(http.StatusOK, res)
}

// PostContest POST /contests
func (h *ContestHandler) PostContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PostContestRequest{}
	err := c.BindAndValidate(&req)
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
		return convertError(err)
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
	ContestID   uuid.UUID       `param:"contestID" validate:"is-uuid"`
	Name        optional.String `json:"name"`
	Link        optional.String `json:"link"`
	Description optional.String `json:"description"`
	Duration    OptionalDuration
}

// PatchContest PATCH /contests/:contestID
func (h *ContestHandler) PatchContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PatchContestRequest{}
	err := c.BindAndValidate(&req)
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

	err = h.srv.UpdateContest(ctx, req.ContestID, &patchReq)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteContest DELETE /contests/:contestID
func (h *ContestHandler) DeleteContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := contestParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}
	err := h.srv.DeleteContest(ctx, req.ContestID)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetContestTeams GET /contests/:contestID/teams
func (h *ContestHandler) GetContestTeams(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := contestParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}
	contestTeams, err := h.srv.GetContestTeams(ctx, req.ContestID)
	if err != nil {
		return convertError(err)
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
	req := teamParams{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}
	contestTeam, err := h.srv.GetContestTeam(ctx, req.ContestID, req.TeamID)
	if err != nil {
		return convertError(err)
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

type PostContestTeamRequest struct {
	ContestID   uuid.UUID `param:"contestID" validate:"is-uuid"`
	Name        string    `json:"name"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Result      string    `json:"result"`
}

type PostContestTeamResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Result string    `json:"result"`
}

// PostContestTeam POST /contests/:contestID/teams
func (h *ContestHandler) PostContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PostContestTeamRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.CreateContestTeamArgs{
		Name:        req.Name,
		Result:      req.Result,
		Link:        req.Link,
		Description: req.Description,
	}
	contestTeam, err := h.srv.CreateContestTeam(ctx, req.ContestID, &args)
	if err != nil {
		return convertError(err)
	}

	res := &PostContestTeamResponse{
		ID:     contestTeam.ID,
		Name:   contestTeam.Name,
		Result: contestTeam.Result,
	}
	return c.JSON(http.StatusCreated, res)
}

type PatchContestTeamRequest struct {
	ContestID   uuid.UUID       `param:"contestID" validate:"is-uuid"`
	TeamID      uuid.UUID       `param:"teamID" validate:"is-uuid"`
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
	req := PatchContestTeamRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.UpdateContestTeamArgs{
		Name:        req.Name,
		Result:      req.Result,
		Link:        req.Link,
		Description: req.Description,
	}

	err = h.srv.UpdateContestTeam(ctx, req.TeamID, &args)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusCreated)
}

// GetContestTeamMember GET /contests/{contestId}/teams/{teamId}/members
func (h *ContestHandler) GetContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := teamParams{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	users, err := h.srv.GetContestTeamMembers(ctx, req.ContestID, req.TeamID)
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
	ContestID uuid.UUID   `param:"contestID" validate:"is-uuid"`
	TeamID    uuid.UUID   `param:"teamID" validate:"is-uuid"`
	Members   []uuid.UUID `json:"members" validate:"required"`
}

// PostContestTeamMember POST /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) PostContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	req := PostContestTeamMember{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.srv.AddContestTeamMembers(ctx, req.TeamID, req.Members)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteContestTeamMember DELETE /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) DeleteContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	req := PostContestTeamMember{} // TODO: 構造体分けたいかも
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.srv.DeleteContestTeamMembers(ctx, req.TeamID, req.Members)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
