package handler

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/optional"
)

//TODO 何月？
var (
	semesterToMonth = [2]time.Month{time.August, time.December}
)

type projectParams struct {
	ProjectID uuid.UUID `param:"projectID" validate:"is-uuid"`
}

// ProjectResponse Portfolioのレスポンスで使うイベント情報
type ProjectResponse struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	Duration domain.ProjectDuration `json:"duration"`
}

type ProjectDetailResponse struct {
	ID          uuid.UUID                      `json:"id"`
	Name        string                         `json:"name"`
	Duration    domain.ProjectDuration         `json:"duration"`
	Link        string                         `json:"link"`
	Description string                         `json:"description"`
	Members     []*ProjectMemberDetailResponse `json:"members"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
}

type ProjectMemberResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"real_name"`
}

type ProjectMemberDetailResponse struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	RealName string                 `json:"real_name"`
	Duration domain.ProjectDuration `json:"duration"`
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
	res := make([]*ProjectResponse, 0, len(projects))
	for _, v := range projects {
		p := &ProjectResponse{
			ID:       v.ID,
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
	req := projectParams{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	project, err := h.service.GetProject(ctx, req.ProjectID)
	if err != nil {
		return convertError(err)
	}
	res := &ProjectDetailResponse{
		ID:          project.ID,
		Name:        project.Name,
		Duration:    convertToProjectDuration(project.Since, project.Until),
		Link:        project.Link,
		Description: project.Description,
		Members:     convertToProjectMembersDetail(project.Members),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)
}

type PostProjectRequest struct {
	Name        string                 `json:"name"`
	Link        string                 `json:"link"`
	Description string                 `json:"description"`
	Duration    domain.ProjectDuration `json:"duration"`
}

// PostProject POST /projects
func (h *ProjectHandler) PostProject(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PostProjectRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	since := semToTime(req.Duration.Since)
	until := semToTime(req.Duration.Until)
	if since.After(until) {
		return convertError(repository.ErrInvalidArg)
	}
	createReq := repository.CreateProjectArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Since:       since,
		Until:       until,
	}
	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return convertError(err)
	}
	res := ProjectResponse{
		ID:       project.ID,
		Name:     project.Name,
		Duration: convertToProjectDuration(project.Since, project.Until),
	}
	return c.JSON(http.StatusCreated, res)
}

type PatchProjectRequest struct {
	ProjectID   uuid.UUID       `param:"projectID" validate:"is-uuid"`
	Name        optional.String `json:"name"`
	Link        optional.String `json:"link"`
	Description optional.String `json:"description"`
	Duration    OptionalProjectDuration
}

func (h *ProjectHandler) PatchProject(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PatchProjectRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	since := optionalSemToTime(req.Duration.Since)
	until := optionalSemToTime(req.Duration.Until)
	if since.Valid && until.Valid && since.Time.After(until.Time) {
		return convertError(repository.ErrInvalidArg)
	}
	patchReq := repository.UpdateProjectArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
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
	req := projectParams{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	members, err := h.service.GetProjectMembers(ctx, req.ProjectID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*ProjectMemberResponse, 0, len(members))
	for _, v := range members {
		m := &ProjectMemberResponse{
			ID:       v.ID,
			Name:     v.Name,
			RealName: v.RealName,
		}
		res = append(res, m)
	}
	return c.JSON(http.StatusOK, res)
}

type AddProjectMembersRequest struct {
	ProjectID uuid.UUID `param:"projectID" validate:"is-uuid"`
	Members   []*MemberIDWithProjectDuration
}

type MemberIDWithProjectDuration struct {
	UserID   uuid.UUID              `json:"userID" validate:"is-uuid"` // TODO: validateしてくれない
	Duration domain.ProjectDuration `json:"duration"`
}

// AddProjectMembers POST /projects/:projectID/members
func (h *ProjectHandler) AddProjectMembers(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := AddProjectMembersRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}
	createReq := make([]*repository.CreateProjectMemberArgs, 0, len(req.Members))
	for _, v := range req.Members {
		if v.UserID == uuid.Nil { // TODO
			return convertError(repository.ErrInvalidArg)
		}
		m := &repository.CreateProjectMemberArgs{
			UserID: v.UserID,
			Since:  semToTime(v.Duration.Since),
			Until:  semToTime(v.Duration.Until),
		}
		createReq = append(createReq, m)
	}
	err = h.service.AddProjectMembers(ctx, req.ProjectID, createReq)
	if err != nil {
		return convertError(err)
	}
	return nil
}

type PutProjectMembersRequest struct {
	ProjectID uuid.UUID   `param:"projectID" validate:"is-uuid"`
	Members   []uuid.UUID `json:"members"`
}

// DeleteProjectMembers DELETE /projects/:projectID/members
func (h *ProjectHandler) DeleteProjectMembers(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := PutProjectMembersRequest{}
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

func convertToProjectMembersDetail(members []*domain.ProjectMember) []*ProjectMemberDetailResponse {
	res := make([]*ProjectMemberDetailResponse, 0, len(members))
	for _, v := range members {
		m := &ProjectMemberDetailResponse{
			ID:       v.UserID,
			Name:     v.Name,
			RealName: v.RealName,
			Duration: domain.ProjectDuration{
				Since: timeToSem(v.Since),
				Until: timeToSem(v.Until),
			},
		}
		res = append(res, m)
	}
	return res
}
