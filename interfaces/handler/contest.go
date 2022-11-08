package handler

import (
	"net/http"
	"time"

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

// GetContests GET /contests
func (h *ContestHandler) GetContests(_c echo.Context) error {
	c := _c.(*Context)

	ctx := c.Request().Context()
	contests, err := h.srv.GetContests(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Contest, len(contests))
	for i, v := range contests {
		res[i] = newContest(v.ID, v.Name, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

// GetContest GET /contests/:contestID
func (h *ContestHandler) GetContest(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	contest, err := h.srv.GetContest(ctx, contestID)
	if err != nil {
		return convertError(err)
	}

	teams := make([]ContestTeam, len(contest.ContestTeams))
	for i, v := range contest.ContestTeams {
		teams[i] = newContestTeam(v.ID, v.Name, v.Result)
	}

	res := newContestDetail(
		newContest(contest.ID, contest.Name, contest.TimeStart, contest.TimeEnd),
		contest.Link,
		contest.Description,
		teams,
	)

	return c.JSON(http.StatusOK, res)
}

// CreateContest POST /contests
func (h *ContestHandler) CreateContest(_c echo.Context) error {
	c := _c.(*Context)

	req := CreateContestJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	createReq := repository.CreateContestArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        optional.StringFrom(req.Link),
		Since:       req.Duration.Since,
		Until:       optional.TimeFrom(req.Duration.Until),
	}

	ctx := c.Request().Context()
	contest, err := h.srv.CreateContest(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}

	contestTeams := make([]ContestTeam, 0, len(contest.ContestTeams))
	for _, team := range contest.ContestTeams {
		contestTeams = append(contestTeams, newContestTeam(team.ID, team.Name, team.Result))
	}
	res := newContestDetail(newContest(contest.ID, contest.Name, contest.TimeStart, contest.TimeEnd), contest.Link, contest.Description, contestTeams)

	return c.JSON(http.StatusCreated, res)
}

// EditContest PATCH /contests/:contestID
func (h *ContestHandler) EditContest(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	req := EditContestJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	patchReq := repository.UpdateContestArgs{
		Name:        optional.StringFrom(req.Name),
		Description: optional.StringFrom(req.Description),
		Link:        optional.StringFrom(req.Link),
	}
	if req.Duration != nil {
		patchReq.Since = optional.TimeFrom(&req.Duration.Since)
		patchReq.Until = optional.TimeFrom(req.Duration.Until)
	}

	ctx := c.Request().Context()
	err = h.srv.UpdateContest(ctx, contestID, &patchReq)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteContest DELETE /contests/:contestID
func (h *ContestHandler) DeleteContest(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	if err := h.srv.DeleteContest(ctx, contestID); err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetContestTeams GET /contests/:contestID/teams
func (h *ContestHandler) GetContestTeams(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	contestTeams, err := h.srv.GetContestTeams(ctx, contestID)
	if err != nil {
		return convertError(err)
	}

	res := make([]ContestTeam, len(contestTeams))
	for i, v := range contestTeams {
		res[i] = newContestTeam(v.ID, v.Name, v.Result)
	}

	return c.JSON(http.StatusOK, res)
}

// GetContestTeams GET /contests/:contestID/teams/:teamID
func (h *ContestHandler) GetContestTeam(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	teamID, err := c.getID(keyContestTeamID)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	contestTeam, err := h.srv.GetContestTeam(ctx, contestID, teamID)
	if err != nil {
		return convertError(err)
	}

	members := make([]User, len(contestTeam.Members))
	for i, v := range contestTeam.Members {
		members[i] = newUser(v.ID, v.Name, v.RealName())
	}

	res := newContestTeamDetail(
		newContestTeam(contestTeam.ID, contestTeam.Name, contestTeam.Result),
		contestTeam.Link,
		contestTeam.Description,
		members,
	)

	return c.JSON(http.StatusOK, res)
}

// AddContestTeam POST /contests/:contestID/teams
func (h *ContestHandler) AddContestTeam(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	req := AddContestTeamJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	args := repository.CreateContestTeamArgs{
		Name:        req.Name,
		Result:      optional.StringFrom(req.Result),
		Link:        optional.StringFrom(req.Link),
		Description: req.Description,
	}

	ctx := c.Request().Context()
	contestTeam, err := h.srv.CreateContestTeam(ctx, contestID, &args)
	if err != nil {
		return convertError(err)
	}

	res := newContestTeam(contestTeam.ID, contestTeam.Name, contestTeam.Result)

	return c.JSON(http.StatusCreated, res)
}

// EditContestTeam PATCH /contests/:contestID/teams/:teamID
func (h *ContestHandler) EditContestTeam(_c echo.Context) error {
	c := _c.(*Context)

	// TODO: contestIDをUpdateContestTeamの引数に含める
	_, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	teamID, err := c.getID(keyContestTeamID)
	if err != nil {
		return convertError(err)
	}

	req := EditContestTeamJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	args := repository.UpdateContestTeamArgs{
		Name:        optional.StringFrom(req.Name),
		Result:      optional.StringFrom(req.Result),
		Link:        optional.StringFrom(req.Link),
		Description: optional.StringFrom(req.Description),
	}

	ctx := c.Request().Context()
	if err = h.srv.UpdateContestTeam(ctx, teamID, &args); err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteContestTeam DELETE /contests/:contestID/teams/:teamID
func (h *ContestHandler) DeleteContestTeam(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	teamID, err := c.getID(keyContestTeamID)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	if err = h.srv.DeleteContestTeam(ctx, contestID, teamID); err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetContestTeamMembers GET /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) GetContestTeamMembers(_c echo.Context) error {
	c := _c.(*Context)

	contestID, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	teamID, err := c.getID(keyContestTeamID)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	users, err := h.srv.GetContestTeamMembers(ctx, contestID, teamID)
	if err != nil {
		return convertError(err)
	}

	res := make([]*User, 0, len(users))
	for _, v := range users {
		res = append(res, &User{
			Id:       v.ID,
			Name:     v.Name,
			RealName: v.RealName(),
		})
	}
	return c.JSON(http.StatusOK, res)
}

// AddContestTeamMembers POST /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) AddContestTeamMembers(_c echo.Context) error {
	c := _c.(*Context)
	ctx := c.Request().Context()

	// TODO: contestIDをAddContestTeamMembersの引数に含める
	_, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	teamID, err := c.getID(keyContestTeamID)
	if err != nil {
		return convertError(err)
	}

	req := MemberIDs{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	err = h.srv.AddContestTeamMembers(ctx, teamID, req.Members)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// EditContestTeamMembers PUT /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) EditContestTeamMembers(_c echo.Context) error {
	c := _c.(*Context)

	// TODO: contestIDをDeleteContestTeamMembersの引数に含める
	_, err := c.getID(keyContestID)
	if err != nil {
		return convertError(err)
	}

	teamID, err := c.getID(keyContestTeamID)
	if err != nil {
		return convertError(err)
	}

	req := MemberIDs{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	if err = h.srv.EditContestTeamMembers(ctx, teamID, req.Members); err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func newContest(id uuid.UUID, name string, since time.Time, until time.Time) Contest {
	return Contest{
		Id:   id,
		Name: name,
		Duration: Duration{
			Since: since,
			Until: &until,
		},
	}
}

func newContestDetail(contest Contest, link string, description string, teams []ContestTeam) ContestDetail {
	return ContestDetail{
		Description: description,
		Duration:    contest.Duration,
		Id:          contest.Id,
		Link:        link,
		Name:        contest.Name,
		Teams:       teams,
	}
}

func newContestTeam(id uuid.UUID, name string, result string) ContestTeam {
	return ContestTeam{
		Id:     id,
		Name:   name,
		Result: result,
	}
}

func newContestTeamDetail(team ContestTeam, link string, description string, members []User) ContestTeamDetail {
	return ContestTeamDetail{
		Description: description,
		Id:          team.Id,
		Link:        link,
		Members:     members,
		Name:        team.Name,
		Result:      team.Result,
	}
}
