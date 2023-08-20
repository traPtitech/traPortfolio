package handler

import (
	"net/http"
	"time"

	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
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
func (h *ContestHandler) GetContests(c echo.Context) error {
	ctx := c.Request().Context()
	contests, err := h.srv.GetContests(ctx)
	if err != nil {
		return err
	}

	res := make([]schema.Contest, len(contests))
	for i, v := range contests {
		res[i] = newContest(v.ID, v.Name, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

// GetContest GET /contests/:contestID
func (h *ContestHandler) GetContest(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	contest, err := h.srv.GetContest(ctx, contestID)
	if err != nil {
		return err
	}

	teams := make([]schema.ContestTeam, len(contest.ContestTeams))
	for i, v := range contest.ContestTeams {
		members := make([]User, len(v.Members))
		for j, ct := range v.Members {
			members[j] = newUser(ct.ID, ct.Name, ct.RealName())
		}
		teams[i] = newContestTeam(v.ID, v.Name, v.Result, members)
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
func (h *ContestHandler) CreateContest(c echo.Context) error {
	req := schema.CreateContestJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	createReq := repository.CreateContestArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        optional.FromPtr(req.Link),
		Since:       req.Duration.Since,
		Until:       optional.FromPtr(req.Duration.Until),
	}

	ctx := c.Request().Context()
	contest, err := h.srv.CreateContest(ctx, &createReq)
	if err != nil {
		return err
	}

	contestTeams := make([]schema.ContestTeam, 0, len(contest.ContestTeams))
	for _, team := range contest.ContestTeams {
		contestTeams = append(contestTeams, newContestTeam(team.ID, team.Name, team.Result, []User{}))
	}
	res := newContestDetail(newContest(contest.ID, contest.Name, contest.TimeStart, contest.TimeEnd), contest.Link, contest.Description, contestTeams)

	return c.JSON(http.StatusCreated, res)
}

// EditContest PATCH /contests/:contestID
func (h *ContestHandler) EditContest(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	req := schema.EditContestJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	patchReq := repository.UpdateContestArgs{
		Name:        optional.FromPtr(req.Name),
		Description: optional.FromPtr(req.Description),
		Link:        optional.FromPtr(req.Link),
	}
	if req.Duration != nil {
		patchReq.Since = optional.FromPtr(&req.Duration.Since)
		patchReq.Until = optional.FromPtr(req.Duration.Until)
	}

	ctx := c.Request().Context()
	err = h.srv.UpdateContest(ctx, contestID, &patchReq)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteContest DELETE /contests/:contestID
func (h *ContestHandler) DeleteContest(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.srv.DeleteContest(ctx, contestID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetContestTeams GET /contests/:contestID/teams
func (h *ContestHandler) GetContestTeams(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	contestTeams, err := h.srv.GetContestTeams(ctx, contestID)
	if err != nil {
		return err
	}

	res := make([]schema.ContestTeam, len(contestTeams))
	for i, v := range contestTeams {
		members := make([]User, len(v.Members))
		for j, ct := range v.Members {
			members[j] = newUser(ct.ID, ct.Name, ct.RealName())
		}
		res[i] = newContestTeam(v.ID, v.Name, v.Result, members)
	}

	return c.JSON(http.StatusOK, res)
}

// GetContestTeams GET /contests/:contestID/teams/:teamID
func (h *ContestHandler) GetContestTeam(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	teamID, err := getID(c, keyContestTeamID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	contestTeam, err := h.srv.GetContestTeam(ctx, contestID, teamID)
	if err != nil {
		return err
	}

	members := make([]schema.User, len(contestTeam.Members))
	for i, v := range contestTeam.Members {
		members[i] = newUser(v.ID, v.Name, v.RealName())
	}

	res := newContestTeamDetail(
		newContestTeam(contestTeam.ID, contestTeam.Name, contestTeam.Result, members),
		contestTeam.Link,
		contestTeam.Description,
	)

	return c.JSON(http.StatusOK, res)
}

// AddContestTeam POST /contests/:contestID/teams
func (h *ContestHandler) AddContestTeam(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	req := schema.AddContestTeamJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	args := repository.CreateContestTeamArgs{
		Name:        req.Name,
		Result:      optional.FromPtr(req.Result),
		Link:        optional.FromPtr(req.Link),
		Description: req.Description,
	}

	ctx := c.Request().Context()
	contestTeam, err := h.srv.CreateContestTeam(ctx, contestID, &args)
	if err != nil {
		return err
	}

	res := newContestTeam(contestTeam.ID, contestTeam.Name, contestTeam.Result, []User{})

	return c.JSON(http.StatusCreated, res)
}

// EditContestTeam PATCH /contests/:contestID/teams/:teamID
func (h *ContestHandler) EditContestTeam(c echo.Context) error {
	// TODO: contestIDをUpdateContestTeamの引数に含める
	_, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	teamID, err := getID(c, keyContestTeamID)
	if err != nil {
		return err
	}

	req := schema.EditContestTeamJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	args := repository.UpdateContestTeamArgs{
		Name:        optional.FromPtr(req.Name),
		Result:      optional.FromPtr(req.Result),
		Link:        optional.FromPtr(req.Link),
		Description: optional.FromPtr(req.Description),
	}

	ctx := c.Request().Context()
	if err = h.srv.UpdateContestTeam(ctx, teamID, &args); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteContestTeam DELETE /contests/:contestID/teams/:teamID
func (h *ContestHandler) DeleteContestTeam(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	teamID, err := getID(c, keyContestTeamID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err = h.srv.DeleteContestTeam(ctx, contestID, teamID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetContestTeamMembers GET /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) GetContestTeamMembers(c echo.Context) error {
	contestID, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	teamID, err := getID(c, keyContestTeamID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	users, err := h.srv.GetContestTeamMembers(ctx, contestID, teamID)
	if err != nil {
		return err
	}

	res := make([]*schema.User, 0, len(users))
	for _, v := range users {
		res = append(res, &schema.User{
			Id:       v.ID,
			Name:     v.Name,
			RealName: v.RealName(),
		})
	}
	return c.JSON(http.StatusOK, res)
}

// AddContestTeamMembers POST /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) AddContestTeamMembers(c echo.Context) error {
	ctx := c.Request().Context()

	// TODO: contestIDをAddContestTeamMembersの引数に含める
	_, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	teamID, err := getID(c, keyContestTeamID)
	if err != nil {
		return err
	}

	req := schema.AddContestTeamMembersJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	err = h.srv.AddContestTeamMembers(ctx, teamID, req.Members)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// EditContestTeamMembers PUT /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) EditContestTeamMembers(c echo.Context) error {
	// TODO: contestIDをDeleteContestTeamMembersの引数に含める
	_, err := getID(c, keyContestID)
	if err != nil {
		return err
	}

	teamID, err := getID(c, keyContestTeamID)
	if err != nil {
		return err
	}

	req := schema.EditContestTeamMembersJSONRequestBody{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err = h.srv.EditContestTeamMembers(ctx, teamID, req.Members); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func newContest(id uuid.UUID, name string, since time.Time, until time.Time) schema.Contest {
	return schema.Contest{
		Id:   id,
		Name: name,
		Duration: schema.Duration{
			Since: since,
			Until: &until,
		},
	}
}

func newContestDetail(contest schema.Contest, link string, description string, teams []schema.ContestTeam) schema.ContestDetail {
	return schema.ContestDetail{
		Description: description,
		Duration:    contest.Duration,
		Id:          contest.Id,
		Link:        link,
		Name:        contest.Name,
		Teams:       teams,
	}
}

func newContestTeam(id uuid.UUID, name string, result string, members []User) schema.ContestTeam {
	return schema.ContestTeam{
		Id:      id,
		Name:    name,
		Result:  result,
		Members: members,
	}
}

func newContestTeamWithoutMembers(id uuid.UUID, name string, result string) schema.ContestTeamWithoutMembers {
	return schema.ContestTeamWithoutMembers{
		Id:     id,
		Name:   name,
		Result: result,
	}
}

func newContestTeamDetail(team schema.ContestTeam, link string, description string) schema.ContestTeamDetail {
	return schema.ContestTeamDetail{
		Description: description,
		Id:          team.Id,
		Link:        link,
		Members:     team.Members,
		Name:        team.Name,
		Result:      team.Result,
	}
}
