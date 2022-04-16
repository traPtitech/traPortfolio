package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/optional"
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
	c := _c.(*Context)
	ctx := c.Request().Context()
	projects, err := h.service.GetProjects(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Project, len(projects))
	for i, v := range projects {
		res[i] = newProject(v.ID, v.Name, convertDuration(v.Duration))
	}

	return c.JSON(http.StatusOK, res)
}

// GetByID GET /projects/:projectID
func (h *ProjectHandler) GetByID(_c echo.Context) error {
	c := _c.(*Context)
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
			convertDuration(v.Duration),
		)
	}

	return c.JSON(http.StatusOK, newProjectDetail(
		newProject(project.ID, project.Name, convertDuration(project.Duration)),
		project.Description,
		project.Link,
		members,
	))
}

// PostProject POST /projects
func (h *ProjectHandler) PostProject(_c echo.Context) error {
	c := _c.(*Context)
	ctx := c.Request().Context()
	req := CreateProjectJSONRequestBody{}
	err := c.BindAndValidate(&req)
	if err != nil {
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

	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusCreated, newProject(
		project.ID,
		project.Name,
		convertDuration(project.Duration),
	))
}

func (h *ProjectHandler) PatchProject(_c echo.Context) error {
	c := _c.(*Context)
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

	err = h.service.UpdateProject(ctx, req.ProjectID, &patchReq)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetProjectMembers GET /projects/:projectID/members
func (h *ProjectHandler) GetProjectMembers(_c echo.Context) error {
	c := _c.(*Context)
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
			YearWithSemesterDuration{}, // TODO: 追加する
		)
	}

	return c.JSON(http.StatusOK, res)
}

// AddProjectMembers POST /projects/:projectID/members
func (h *ProjectHandler) AddProjectMembers(_c echo.Context) error {
	c := _c.(*Context)
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
	err = h.service.AddProjectMembers(ctx, req.ProjectID, createReq)
	if err != nil {
		return convertError(err)
	}
	return nil
}

// DeleteProjectMembers DELETE /projects/:projectID/members
func (h *ProjectHandler) DeleteProjectMembers(_c echo.Context) error {
	c := _c.(*Context)
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

func newProject(id uuid.UUID, name string, duration YearWithSemesterDuration) Project {
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

func newProjectMember(user User, duration YearWithSemesterDuration) ProjectMember {
	return ProjectMember{
		User:     user,
		Duration: duration,
	}
}
