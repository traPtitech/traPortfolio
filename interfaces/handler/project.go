package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type ProjectHandler struct {
	project repository.ProjectRepository
}

func NewProjectHandler(r repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{r}
}

// GetProjects GET /projects
func (h *ProjectHandler) GetProjects(c echo.Context) error {
	ctx := c.Request().Context()
	projects, err := h.project.GetProjects(ctx)
	if err != nil {
		return err
	}

	res := make([]schema.Project, len(projects))
	for i, v := range projects {
		res[i] = newProject(v.ID, v.Name, schema.ConvertDuration(v.Duration))
	}

	return c.JSON(http.StatusOK, res)
}

// GetProject GET /projects/:projectID
func (h *ProjectHandler) GetProject(c echo.Context) error {
	projectID, err := getID(c, keyProject)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := h.project.GetProject(ctx, projectID)
	if err != nil {
		return err
	}

	members := make([]schema.ProjectMember, len(project.Members))
	for i, v := range project.Members {
		members[i] = newProjectMember(
			newUser(v.User.ID, v.User.Name, v.User.RealName()),
			schema.ConvertDuration(v.Duration),
		)
	}

	return c.JSON(http.StatusOK, newProjectDetail(
		newProject(project.ID, project.Name, schema.ConvertDuration(project.Duration)),
		project.Description,
		project.Link,
		members,
	))
}

// CreateProject POST /projects
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	req := schema.CreateProjectRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	createReq := repository.CreateProjectArgs{
		Name:          req.Name,
		Description:   req.Description,
		Link:          optional.FromPtr(req.Link),
		SinceYear:     req.Duration.Since.Year,
		SinceSemester: int(req.Duration.Since.Semester),
	}

	if req.Duration.Until != nil {
		createReq.UntilYear = req.Duration.Until.Year
		createReq.UntilSemester = int(req.Duration.Until.Semester)
	}
	d := domain.NewYearWithSemesterDuration(createReq.SinceYear, createReq.SinceSemester, createReq.UntilYear, createReq.UntilSemester)
	if !d.IsValid() {
		return repository.ErrInvalidArg
	}

	ctx := c.Request().Context()
	project, err := h.project.CreateProject(ctx, &createReq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, newProject(
		project.ID,
		project.Name,
		schema.ConvertDuration(project.Duration),
	))
}

// EditProject PATCH /projects/:projectID
func (h *ProjectHandler) EditProject(c echo.Context) error {
	projectID, err := getID(c, keyProject)
	if err != nil {
		return err
	}

	req := schema.EditProjectRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	patchReq := repository.UpdateProjectArgs{
		Name:        optional.FromPtr(req.Name),
		Description: optional.FromPtr(req.Description),
		Link:        optional.FromPtr(req.Link),
	}

	if d := req.Duration; d != nil {
		sinceYear := int64(d.Since.Year)
		sinceSemester := int64(d.Since.Semester)
		patchReq.SinceYear = optional.FromPtr(&sinceYear)
		patchReq.SinceSemester = optional.FromPtr(&sinceSemester)

		if d.Until != nil {
			untilYear := int64(d.Until.Year)
			untilSemester := int64(d.Until.Semester)
			patchReq.UntilYear = optional.FromPtr(&untilYear)
			patchReq.UntilSemester = optional.FromPtr(&untilSemester)
		}
	}

	ctx := c.Request().Context()
	{
		old, err := h.project.GetProject(ctx, projectID)
		if err != nil {
			return err
		}

		d := old.Duration
		if sy, ok := patchReq.SinceYear.V(); ok {
			if ss, ok := patchReq.SinceSemester.V(); ok {
				d.Since.Year = int(sy)
				d.Since.Semester = int(ss)
			}
		}
		if uy, ok := patchReq.UntilYear.V(); ok {
			if us, ok := patchReq.UntilSemester.V(); ok {
				d.Until.Year = int(uy)
				d.Until.Semester = int(us)
			}
		}

		if !d.IsValid() {
			return repository.ErrInvalidArg
		}
	}

	err = h.project.UpdateProject(ctx, projectID, &patchReq)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetProjectMembers GET /projects/:projectID/members
func (h *ProjectHandler) GetProjectMembers(c echo.Context) error {
	projectID, err := getID(c, keyProject)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	members, err := h.project.GetProjectMembers(ctx, projectID)
	if err != nil {
		return err
	}

	res := make([]schema.ProjectMember, len(members))
	for i, v := range members {
		res[i] = newProjectMember(
			newUser(v.User.ID, v.User.Name, v.User.RealName()),
			schema.ConvertDuration(v.Duration),
		)
	}

	return c.JSON(http.StatusOK, res)
}

// EditProjectMembers POST /projects/:projectID/members
func (h *ProjectHandler) EditProjectMembers(c echo.Context) error {
	projectID, err := getID(c, keyProject)
	if err != nil {
		return err
	}

	req := schema.EditProjectMembersRequest{}
	if err := c.Bind(&req); err != nil {
		return err
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

		d := domain.NewYearWithSemesterDuration(m.SinceYear, m.SinceSemester, m.UntilYear, m.UntilSemester)
		if !d.IsValid() {
			return repository.ErrInvalidArg
		}

		createReq = append(createReq, m)
	}

	ctx := c.Request().Context()
	err = h.project.EditProjectMembers(ctx, projectID, createReq)
	if err != nil {
		return err
	}

	return nil
}

func newProject(id uuid.UUID, name string, duration schema.YearWithSemesterDuration) schema.Project {
	return schema.Project{
		Id:       id,
		Name:     name,
		Duration: duration,
	}
}

func newProjectDetail(project schema.Project, description string, link string, members []schema.ProjectMember) schema.ProjectDetail {
	return schema.ProjectDetail{
		Description: description,
		Duration:    project.Duration,
		Link:        link,
		Id:          project.Id,
		Members:     members,
		Name:        project.Name,
	}
}

func newProjectMember(user schema.User, duration schema.YearWithSemesterDuration) schema.ProjectMember {
	return schema.ProjectMember{
		Duration: duration,
		Id:       user.Id,
		Name:     user.Name,
		RealName: user.RealName,
	}
}
