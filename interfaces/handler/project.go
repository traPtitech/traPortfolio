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
	res := make([]*Project, 0, len(projects))
	for _, v := range projects {
		p := &Project{
			Id:       v.ID,
			Name:     v.Name,
			Duration: convertToProjectDuration(v.Since, v.Until),
		}
		res = append(res, p)
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
		until := timeToSem(v.Until)
		members[i] = ProjectMember{
			User: User{
				Id:       v.UserID,
				Name:     v.Name,
				RealName: &v.RealName,
			},
			Duration: []ProjectDuration{
				{
					Since: timeToSem(v.Since),
					Until: &until,
				},
			},
		}
	}

	res := &ProjectDetail{
		Project: Project{
			Id:       project.ID,
			Name:     project.Name,
			Duration: convertToProjectDuration(project.Since, project.Until),
		},
		Link:        &project.Link,
		Description: project.Description,
		Members:     members,
		CreatedAt:   &project.CreatedAt,
		UpdatedAt:   &project.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)
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

	since := semToTime(req.Duration.Since)
	until := semToTime(*req.Duration.Until)
	if since.After(until) {
		return convertError(repository.ErrInvalidArg)
	}
	createReq := repository.CreateProjectArgs{
		Name:        *req.Name,
		Description: *req.Description,
		Link:        *req.Link,
		Since:       since,
		Until:       until,
	}
	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}
	res := Project{
		Id:       project.ID,
		Name:     project.Name,
		Duration: convertToProjectDuration(project.Since, project.Until),
	}
	return c.JSON(http.StatusCreated, res)
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

	since := optionalSemToTime(req.Duration.Since)
	until := optionalSemToTime(*req.Duration.Until)
	if since.Valid && until.Valid && since.Time.After(until.Time) {
		return convertError(repository.ErrInvalidArg)
	}
	patchReq := repository.UpdateProjectArgs{
		Name:        optional.StringFrom(*req.Name),
		Description: optional.StringFrom(*req.Description),
		Link:        optional.StringFrom(*req.Link),
		Since:       since,
		Until:       until,
	}

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
	res := make([]*ProjectMember, 0, len(members))
	for _, v := range members {
		m := &ProjectMember{
			User: User{
				Id:       v.ID,
				Name:     v.Name,
				RealName: &v.RealName,
			},
			// Duration: , //TODO
		}
		res = append(res, m)
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
			Since:  semToTime(v.Duration.Since),
			Until:  semToTime(*v.Duration.Until),
		}
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
