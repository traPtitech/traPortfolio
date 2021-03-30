package handler

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/project_service"
)

//TODO 何月？
var (
	semesterToMonth = [2]time.Month{time.August, time.December}
)

type PatchProject struct {
	Name        string                 `json:"name"`
	Link        string                 `json:"link"`
	Description string                 `json:"description"`
	Duration    domain.ProjectDuration `json:"duration"`
}

type ProjectHandler struct {
	service service.ProjectService
}

// PostProjectResponse Portfolioのレスポンスで使うイベント情報
type PostProjectResponse struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	Duration domain.ProjectDuration `json:"duration"`
}

func NewProjectHandler(s service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// PostProject POST /projects
func (h *ProjectHandler) PostProject(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	req := &PatchProject{}
	err := c.Bind(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	createReq := repository.CreateProjectArgs{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Since:       SemToTime(req.Duration.Since),
		Until:       SemToTime(req.Duration.Until),
	}
	project, err := h.service.CreateProject(ctx, &createReq)
	if err != nil {
		return err
	}
	res := PostProjectResponse{
		ID:   project.ID,
		Name: project.Name,
		Duration: domain.ProjectDuration{
			Since: TimeToSem(project.Since),
			Until: TimeToSem(project.Until),
		},
	}
	return c.JSON(http.StatusCreated, res)
}

//TODO 関数名変えたい
func SemToTime(date domain.YearWithSemester) time.Time {
	year := int(date.Year)
	month := semesterToMonth[date.Semester]
	return time.Date(year, month, 1, 0, 0, 0, 0, &time.Location{})
}

func TimeToSem(t time.Time) domain.YearWithSemester {
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
