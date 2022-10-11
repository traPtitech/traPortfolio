package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(s service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// GetProjects GET /projects
func (h *ProjectHandler) GetProjects(_c echo.Context) error {
	c := _c.(*Context)

	ctx := c.Request().Context()
	projects, err := h.service.GetProjects(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Project, len(projects))
	for i, v := range projects {
		res[i] = newProject(v.ID, v.Name, ConvertDuration(v.Duration))
	}

	return c.JSON(http.StatusOK, res)
}

// GetProject GET /projects/:projectID
func (h *ProjectHandler) GetProject(_c echo.Context) error {
	c := _c.(*Context)

	projectID, err := c.getID(keyProject)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	project, err := h.service.GetProject(ctx, projectID)
	if err != nil {
		return convertError(err)
	}

	members := make([]ProjectMember, len(project.Members))
	for i, v := range project.Members {
		members[i] = newProjectMember(
			newUser(v.User.ID, v.User.Name, v.User.RealName),
			ConvertDuration(v.Duration),
		)
	}

	return c.JSON(http.StatusOK, newProjectDetail(
		newProject(project.ID, project.Name, ConvertDuration(project.Duration)),
		project.Description,
		project.Link,
		members,
	))
}

// CreateProject POST /projects
func (h *ProjectHandler) CreateProject(_c echo.Context) error {
	c := _c.(*Context)

	req := CreateProjectJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	createReq := repository.CreateProjectArgs{
		Name:          req.Name,
		Description:   req.Description,
		Link:          optional.StringFrom(req.Link),
		SinceYear:     req.Duration.Since.Year,
		SinceSemester: int(req.Duration.Since.Semester),
	}

	if req.Duration.Until != nil {
		createReq.UntilYear = req.Duration.Until.Year
		createReq.UntilSemester = int(req.Duration.Until.Semester)
	}

	ctx := c.Request().Context()
	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusCreated, newProject(
		project.ID,
		project.Name,
		ConvertDuration(project.Duration),
	))
}

// EditProject PATCH /projects/:projectID
func (h *ProjectHandler) EditProject(_c echo.Context) error {
	c := _c.(*Context)

	projectID, err := c.getID(keyProject)
	if err != nil {
		return convertError(err)
	}

	req := EditProjectJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	patchReq := repository.UpdateProjectArgs{
		Name:        optional.StringFrom(req.Name),
		Description: optional.StringFrom(req.Description),
		Link:        optional.StringFrom(req.Link),
	}

	if d := req.Duration; d != nil {
		sinceYear := int64(d.Since.Year)
		sinceSemester := int64(d.Since.Semester)
		patchReq.SinceYear = optional.Int64From(&sinceYear)
		patchReq.SinceSemester = optional.Int64From(&sinceSemester)

		if d.Until != nil {
			untilYear := int64(d.Until.Year)
			untilSemester := int64(d.Until.Semester)
			patchReq.UntilYear = optional.Int64From(&untilYear)
			patchReq.UntilSemester = optional.Int64From(&untilSemester)
		}
	}

	ctx := c.Request().Context()
	err = h.service.UpdateProject(ctx, projectID, &patchReq)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetProjectMembers GET /projects/:projectID/members
func (h *ProjectHandler) GetProjectMembers(_c echo.Context) error {
	c := _c.(*Context)

	projectID, err := c.getID(keyProject)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	members, err := h.service.GetProjectMembers(ctx, projectID)
	if err != nil {
		return convertError(err)
	}

	res := make([]ProjectMember, len(members))
	for i, v := range members {
		res[i] = newProjectMember(
			newUser(v.User.ID, v.User.Name, v.User.RealName),
			ConvertDuration(v.Duration),
		)
	}

	return c.JSON(http.StatusOK, res)
}

// AddProjectMembers POST /projects/:projectID/members
func (h *ProjectHandler) AddProjectMembers(_c echo.Context) error {
	c := _c.(*Context)

	projectID, err := c.getID(keyProject)
	if err != nil {
		return convertError(err)
	}

	req := AddProjectMembersJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	createReq := make([]*repository.CreateProjectMemberArgs, 0, len(req.Members))
	for _, v := range req.Members {
		m := &repository.CreateProjectMemberArgs{
			UserID:        v.UserId,
			SinceYear:     int(v.Duration.Since.Year),
			SinceSemester: int(v.Duration.Since.Semester),
		}

		if v.Duration.Until != nil {
			m.UntilYear = int(v.Duration.Until.Year)
			m.UntilSemester = int(v.Duration.Until.Semester)
		}

		createReq = append(createReq, m)
	}

	ctx := c.Request().Context()
	err = h.service.AddProjectMembers(ctx, projectID, createReq)
	if err != nil {
		return convertError(err)
	}

	return nil
}

// DeleteProjectMembers DELETE /projects/:projectID/members
func (h *ProjectHandler) DeleteProjectMembers(_c echo.Context) error {
	c := _c.(*Context)

	projectID, err := c.getID(keyProject)
	if err != nil {
		return convertError(err)
	}

	req := DeleteProjectMembersJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	err = h.service.DeleteProjectMembers(ctx, projectID, req.Members)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func newProject(id uuid.UUID, name string, duration YearWithSemesterDuration) Project {
	return Project{
		Id:       id,
		Name:     name,
		Duration: duration,
	}
}

func newProjectDetail(project Project, description string, link string, members []ProjectMember) ProjectDetail {
	return ProjectDetail{
		Description: description,
		Duration:    project.Duration,
		Link:        link,
		Id:          project.Id,
		Members:     members,
		Name:        project.Name,
	}
}

func newProjectMember(user User, duration YearWithSemesterDuration) ProjectMember {
	return ProjectMember{
		Duration: duration,
		Id:       user.Id,
		Name:     user.Name,
		RealName: user.RealName,
	}
}
