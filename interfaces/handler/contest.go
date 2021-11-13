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

type ContestIDInPath struct {
	ContestID uuid.UUID `param:"contestID" validate:"is-uuid"`
}

type TeamIDInPath struct {
	TeamID uuid.UUID `param:"teamID" validate:"is-uuid"`
}

type ContestHandler struct {
	srv service.ContestService
}

// NewContestHandler creates a ContestHandler
func NewContestHandler(service service.ContestService) *ContestHandler {
	return &ContestHandler{service}
}

// GetContests GET /contests
func (h *ContestHandler) GetContests(_c echo.Context) error {
	c := Context{_c}
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

func (h *ContestHandler) GetContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := ContestIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	contest, err := h.srv.GetContest(ctx, req.ContestID)
	if err != nil {
		return convertError(err)
	}

	teams := make([]ContestTeam, len(contest.Teams))
	for i, v := range contest.Teams {
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

// PostContest POST /contests
func (h *ContestHandler) PostContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PostContestJSONRequestBody{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	createReq := repository.CreateContestArgs{
		Name:        req.Name,
		Description: req.Description,
		Since:       req.Duration.Since,
	}
	if req.Link != nil {
		createReq.Link = optional.StringFrom(*req.Link)
	}
	if req.Duration.Until != nil {
		createReq.Until = optional.TimeFrom(*req.Duration.Until)
	}

	contest, err := h.srv.CreateContest(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}

	res := newContest(contest.ID, contest.Name, contest.TimeStart, contest.TimeEnd)

	return c.JSON(http.StatusCreated, res)
}

// PatchContest PATCH /contests/:contestID
func (h *ContestHandler) PatchContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ContestIDInPath
		EditContestJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	patchReq := repository.UpdateContestArgs{}
	if req.Name != nil {
		patchReq.Name = optional.StringFrom(*req.Name)
	}
	if req.Description != nil {
		patchReq.Description = optional.StringFrom(*req.Description)
	}
	if req.Link != nil {
		patchReq.Link = optional.StringFrom(*req.Link)
	}
	if req.Duration != nil {
		patchReq.Since = optional.TimeFrom(req.Duration.Since)
		if req.Duration.Until != nil {
			patchReq.Until = optional.TimeFrom(*req.Duration.Until)
		}
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
	req := ContestIDInPath{}
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
	req := ContestIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	contestTeams, err := h.srv.GetContestTeams(ctx, req.ContestID)
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
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ContestIDInPath
		TeamIDInPath
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}
	contestTeam, err := h.srv.GetContestTeam(ctx, req.ContestID, req.TeamID)
	if err != nil {
		return convertError(err)
	}

	members := make([]User, len(contestTeam.Members))
	for i, v := range contestTeam.Members {
		members[i] = newUser(v.ID, v.Name, v.RealName)
	}

	res := newContestTeamDetail(
		newContestTeam(contestTeam.ID, contestTeam.Name, contestTeam.Result),
		contestTeam.Link,
		contestTeam.Description,
		members,
	)

	return c.JSON(http.StatusOK, res)
}

// PostContestTeam POST /contests/:contestID/teams
func (h *ContestHandler) PostContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ContestIDInPath
		PostContestTeamJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	args := repository.CreateContestTeamArgs{
		Name:        req.Name,
		Description: req.Description,
	}
	if req.Result != nil {
		args.Result = optional.StringFrom(*req.Result)
	}
	if req.Link != nil {
		args.Link = optional.StringFrom(*req.Link)
	}

	contestTeam, err := h.srv.CreateContestTeam(ctx, req.ContestID, &args)
	if err != nil {
		return convertError(err)
	}

	res := newContestTeam(contestTeam.ID, contestTeam.Name, contestTeam.Result)

	return c.JSON(http.StatusCreated, res)
}

// PatchContestTeam PATCH /contests/:contestID
func (h *ContestHandler) PatchContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	req := struct {
		ContestIDInPath
		TeamIDInPath
		EditContestTeamJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.UpdateContestTeamArgs{}
	if req.Name != nil {
		args.Name = optional.StringFrom(*req.Name)
	}
	if req.Result != nil {
		args.Result = optional.StringFrom(*req.Result)
	}
	if req.Link != nil {
		args.Link = optional.StringFrom(*req.Link)
	}
	if req.Description != nil {
		args.Description = optional.StringFrom(*req.Description)
	}

	err = h.srv.UpdateContestTeam(ctx, req.TeamID, &args)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetContestTeamMember GET /contests/{contestId}/teams/{teamId}/members
func (h *ContestHandler) GetContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ContestIDInPath
		TeamIDInPath
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	users, err := h.srv.GetContestTeamMembers(ctx, req.ContestID, req.TeamID)
	if err != nil {
		return convertError(err)
	}

	res := make([]*User, 0, len(users))
	for _, v := range users {
		res = append(res, &User{
			Id:       v.ID,
			Name:     v.Name,
			RealName: v.RealName,
		})
	}
	return c.JSON(http.StatusOK, res)
}

// PostContestTeamMember POST /contests/:contestID/teams/:teamID/members
func (h *ContestHandler) PostContestTeamMember(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	req := struct {
		ContestIDInPath
		TeamIDInPath
		MemberIDs
	}{}
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
	req := struct {
		ContestIDInPath
		TeamIDInPath
		MemberIDs
	}{}
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
		Contest:     contest,
		Link:        link,
		Description: description,
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
		ContestTeam: team,
		Link:        link,
		Description: description,
		Members:     members,
	}
}
