package handler

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/optional"
)

//TODO 何月？
var (
	semesterToMonth = [2]time.Month{time.August, time.December}
)

type ProjectIDInPath struct {
	ProjectID uuid.UUID `param:"projectID" validate:"is-uuid"`
}

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(s service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// GetAll GET /projects
func (h *ProjectHandler) GetAll(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	projects, err := h.service.GetProjects(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Project, len(projects))
	for i, v := range projects {
		res[i] = newProject(v.ID, v.Name, convertToProjectDuration(v.Since, v.Until))
	}

	return c.JSON(http.StatusOK, res)
}

// GetByID GET /projects/:projectID
func (h *ProjectHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
	req := ProjectIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	project, err := h.service.GetProject(ctx, req.ProjectID)
	if err != nil {
		return convertError(err)
	}

	members := make([]ProjectMember, len(project.Members))
	for i, v := range project.Members {
		members[i] = newProjectMember(
			newUser(v.UserID, v.Name, v.RealName),
			[]ProjectDuration{newProjectDuration(timeToSem(v.Since), timeToSem(v.Until))},
		)
	}

	return c.JSON(http.StatusOK, newProjectDetail(
		newProject(project.ID, project.Name, convertToProjectDuration(project.Since, project.Until)),
		project.Description,
		project.Link,
		members,
	))
}

// PostProject POST /projects
func (h *ProjectHandler) PostProject(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PostProjectJSONRequestBody{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	createReq := repository.CreateProjectArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        optional.StringFrom(req.Link),
	}
	since := semToTime(req.Duration.Since)
	if req.Duration.Until != nil {
		until := semToTime(*req.Duration.Until)
		if since.After(until) {
			return convertError(repository.ErrInvalidArg)
		}
		createReq.Until = until
	}
	createReq.Since = since

	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusCreated, newProject(
		project.ID,
		project.Name,
		convertToProjectDuration(project.Since, project.Until),
	))
}

func (h *ProjectHandler) PatchProject(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ProjectIDInPath
		EditProjectJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	patchReq := repository.UpdateProjectArgs{
		Name:        optional.StringFrom(req.Name),
		Description: optional.StringFrom(req.Description),
		Link:        optional.StringFrom(req.Link),
	}
	since := optionalSemToTime(req.Duration.Since)
	if req.Duration.Until != nil {
		until := optionalSemToTime(*req.Duration.Until)
		if since.Valid && until.Valid && since.Time.After(until.Time) {
			return convertError(repository.ErrInvalidArg)
		}
		patchReq.Until = until
	}
	patchReq.Since = since

	err = h.service.UpdateProject(ctx, req.ProjectID, &patchReq)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetProjectMembers GET /projects/:projectID/members
func (h *ProjectHandler) GetProjectMembers(_c echo.Context) error {
	c := Context{_c}
	req := ProjectIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	members, err := h.service.GetProjectMembers(ctx, req.ProjectID)
	if err != nil {
		return convertError(err)
	}

	res := make([]ProjectMember, len(members))
	for i, v := range members {
		res[i] = newProjectMember(
			newUser(v.ID, v.Name, v.RealName),
			[]ProjectDuration{}, // TODO: 追加する
		)
	}

	return c.JSON(http.StatusOK, res)
}

// AddProjectMembers POST /projects/:projectID/members
func (h *ProjectHandler) AddProjectMembers(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ProjectIDInPath
		AddProjectMembersJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	createReq := make([]*repository.CreateProjectMemberArgs, 0, len(req.Members))
	for _, v := range req.Members {
		m := &repository.CreateProjectMemberArgs{
			UserID: v.UserId,
		}
		since := semToTime(v.Duration.Since)
		if v.Duration.Until != nil {
			until := semToTime(*v.Duration.Until)
			if since.After(until) {
				return convertError(repository.ErrInvalidArg)
			}
			m.Until = until
		}
		m.Since = since
		createReq = append(createReq, m)
	}
	err = h.service.AddProjectMembers(ctx, req.ProjectID, createReq)
	if err != nil {
		return convertError(err)
	}
	return nil
}

// DeleteProjectMembers DELETE /projects/:projectID/members
func (h *ProjectHandler) DeleteProjectMembers(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := struct {
		ProjectIDInPath
		DeleteProjectMembersJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	err = h.service.DeleteProjectMembers(ctx, req.ProjectID, req.Members)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func newProject(id uuid.UUID, name string, duration ProjectDuration) Project {
	return Project{
		Id:       id,
		Name:     name,
		Duration: duration,
	}
}

func newProjectDetail(project Project, description string, link string, members []ProjectMember) ProjectDetail {
	return ProjectDetail{
		Project:     project,
		Description: description,
		Link:        link,
		Members:     members,
	}
}

func newProjectMember(user User, duration []ProjectDuration) ProjectMember {
	return ProjectMember{
		User:     user,
		Duration: duration,
	}
}

func newProjectDuration(since YearWithSemester, until YearWithSemester) ProjectDuration {
	return ProjectDuration{
		Since: since,
		Until: &until,
	}
}
