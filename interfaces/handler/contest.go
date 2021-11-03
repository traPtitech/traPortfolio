package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type contestIDInPath struct {
	ContestID uuid.UUID `param:"contestID" validate:"is-uuid"`
}

type teamIDInPath struct {
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
	res := make([]*Contest, 0, len(contests))
	for _, v := range contests {
		res = append(res, &Contest{
			Id:   v.ID,
			Name: v.Name,
			Duration: Duration{
				Since: v.TimeStart,
				Until: &v.TimeEnd,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *ContestHandler) GetContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := contestIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	contest, err := h.srv.GetContest(ctx, req.ContestID)
	if err != nil {
		return convertError(err)
	}

	teams := make([]*ContestTeam, 0, len(contest.Teams))
	for _, v := range contest.Teams {
		teams = append(teams, &ContestTeam{
			Id:     v.ID,
			Name:   v.Name,
			Result: &v.Result,
		})
	}

	res := &ContestDetail{
		Contest: Contest{
			Id:   contest.ID,
			Name: contest.Name,
			Duration: Duration{
				Since: contest.TimeStart,
				Until: &contest.TimeEnd,
			},
		},
		Link:        &contest.Link,
		Description: contest.Description,
		// Teams:       teams, //TODO
	}

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
		Name:        *req.Name,
		Description: *req.Description,
		Link:        *req.Link,
		Since:       req.Duration.Since,
		Until:       *req.Duration.Until,
	}

	contest, err := h.srv.CreateContest(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}
	res := Contest{
		Id:   contest.ID,
		Name: contest.Name,
		Duration: Duration{
			Since: contest.TimeStart,
			Until: &contest.TimeEnd,
		},
	}
	return c.JSON(http.StatusCreated, res)
}

// PatchContest PATCH /contests/:contestID
func (h *ContestHandler) PatchContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		contestIDInPath
		EditContestJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	patchReq := repository.UpdateContestArgs{
		Name:        optional.StringFrom(req.Name),
		Description: optional.StringFrom(req.Description),
		Link:        optional.StringFrom(req.Link),
		Since:       optional.TimeFrom(&req.Duration.Since),
		Until:       optional.TimeFrom(req.Duration.Until),
	}

	err = h.srv.UpdateContest(ctx, req.ContestID, &patchReq)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusCreated)
}

// DeleteContest DELETE /contests/:contestID
func (h *ContestHandler) DeleteContest(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := contestIDInPath{}
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
	req := contestIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}
	contestTeams, err := h.srv.GetContestTeams(ctx, req.ContestID)
	if err != nil {
		return convertError(err)
	}

	res := make([]*ContestTeam, 0, len(contestTeams))
	for _, v := range contestTeams {
		ct := &ContestTeam{
			Id:     v.ID,
			Name:   v.Name,
			Result: &v.Result,
		}
		res = append(res, ct)
	}
	return c.JSON(http.StatusOK, res)
}

// GetContestTeams GET /contests/:contestID/teams/:teamID
func (h *ContestHandler) GetContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		contestIDInPath
		teamIDInPath
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}
	contestTeam, err := h.srv.GetContestTeam(ctx, req.ContestID, req.TeamID)
	if err != nil {
		return convertError(err)
	}

	members := make([]*User, 0, len(contestTeam.Members))
	for _, user := range contestTeam.Members {
		members = append(members, &User{
			Id:       user.ID,
			Name:     user.Name,
			RealName: &user.RealName,
		})
	}

	res := &ContestTeamDetail{
		ContestTeam: ContestTeam{
			Id:     contestTeam.ID,
			Name:   contestTeam.Name,
			Result: &contestTeam.Result,
		},
		Link:        &contestTeam.Link,
		Description: contestTeam.Description,
		// Members:     members, //TODO
	}
	return c.JSON(http.StatusOK, res)
}

// PostContestTeam POST /contests/:contestID/teams
func (h *ContestHandler) PostContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		contestIDInPath
		PostContestTeamJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.CreateContestTeamArgs{
		Name:        *req.Name,
		Result:      *req.Result,
		Link:        *req.Link,
		Description: *req.Description,
	}
	contestTeam, err := h.srv.CreateContestTeam(ctx, req.ContestID, &args)
	if err != nil {
		return convertError(err)
	}

	res := &ContestTeam{
		Id:     contestTeam.ID,
		Name:   contestTeam.Name,
		Result: &contestTeam.Result,
	}
	return c.JSON(http.StatusCreated, res)
}

// PatchContestTeam PATCH /contests/:contestID
func (h *ContestHandler) PatchContestTeam(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	// todo contestIDが必要ない
	req := struct {
		contestIDInPath
		teamIDInPath
		EditContestTeamJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	args := repository.UpdateContestTeamArgs{
		Name:        optional.StringFrom(req.Name),
		Result:      optional.StringFrom(req.Result),
		Link:        optional.StringFrom(req.Link),
		Description: optional.StringFrom(req.Description),
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
	req := struct {
		contestIDInPath
		teamIDInPath
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
			RealName: &v.RealName,
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
		contestIDInPath
		teamIDInPath
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
		contestIDInPath
		teamIDInPath
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
