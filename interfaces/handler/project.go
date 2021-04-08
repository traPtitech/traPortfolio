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

// ProjectResponse Portfolioのレスポンスで使うイベント情報
type ProjectResponse struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	Duration domain.ProjectDuration `json:"duration"`
}

type ProjectDetailResponse struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Duration    domain.ProjectDuration  `json:"duration"`
	Link        string                  `json:"link"`
	Description string                  `json:"description"`
	Members     []*domain.ProjectMember `json:"members"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
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
		return err
	}
	res := make([]*ProjectResponse, 0, len(projects))
	for _, v := range projects {
		res = append(res, &ProjectResponse{
			ID:   v.ID,
			Name: v.Name,
			Duration: domain.ProjectDuration{
				Since: timeToSem(v.Since),
				Until: timeToSem(v.Until),
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

// GetByID GET /projects/:projectID
func (h *ProjectHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("projectID")
	id := uuid.FromStringOrNil(_id)
	project, err := h.service.GetProject(ctx, id)
	if err != nil {
		return err
	}
	res := &ProjectDetailResponse{
		ID:   project.ID,
		Name: project.Name,
		Duration: domain.ProjectDuration{
			Since: timeToSem(project.Since),
			Until: timeToSem(project.Until),
		},
		Link: project.Link,
		// Members:   project.Members, //TODO
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
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
	req := &PostProjectRequest{}
	// todo validation
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	createReq := repository.CreateProjectArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Since:       semToTime(req.Duration.Since),
		Until:       semToTime(req.Duration.Until),
	}
	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return err
	}
	res := ProjectResponse{
		ID:   project.ID,
		Name: project.Name,
		Duration: domain.ProjectDuration{
			Since: timeToSem(project.Since),
			Until: timeToSem(project.Until),
		},
	}
	return c.JSON(http.StatusCreated, res)
}

type PatchProjectRequest struct {
	Name        optional.String `json:"name"`
	Link        optional.String `json:"link"`
	Description optional.String `json:"description"`
	Duration    OptionalProjectDuration
}

func (h *ProjectHandler) PatchProject(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("projectID")
	id := uuid.FromStringOrNil(_id)
	req := &PatchProjectRequest{}
	// todo validation
	err := c.BindAndValidate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	patchReq := repository.UpdateProjectArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Since:       optionalSemToTime(req.Duration.Since),
		Until:       optionalSemToTime(req.Duration.Until),
	}

	err = h.service.UpdateProject(ctx, id, &patchReq)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func semToTime(date domain.YearWithSemester) time.Time {
	year := int(date.Year)
	month := semesterToMonth[date.Semester]
	return time.Date(year, month, 1, 0, 0, 0, 0, &time.Location{})
}

func timeToSem(t time.Time) domain.YearWithSemester {
	year := uint(t.Year())
	var semester uint
	for i, v := range semesterToMonth {
		if v == t.Month() {
			semester = uint(i)
		}
	}
	return domain.YearWithSemester{
		Year:     year,
		Semester: semester,
	}
}

func optionalSemToTime(date OptionalYearWithSemester) optional.Time {
	t := optional.Time{}
	if date.Year.Valid && date.Semester.Valid {
		year := int(date.Year.Int64)
		month := semesterToMonth[date.Semester.Int64]
		t.Time, t.Valid = time.Date(year, month, 1, 0, 0, 0, 0, &time.Location{}), true
	} else {
		t.Valid = false
	}
	return t
}
