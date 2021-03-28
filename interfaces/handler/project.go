package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	service "github.com/traPtitech/traPortfolio/usecases/service/project_service"
)

type PatchProject struct {
	Name        string                 `json:"name"`
	Link        string                 `json:"link"`
	Description string                 `json:"description"`
	Duration    domain.ProjectDuration `json:"duration"`
}

type ProjectHandler struct {
	ProjectService service.ProjectService
}

// ProjectResponse Portfolioのレスポンスで使うイベント情報
type ProjectResponse struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	Duration domain.ProjectDuration `json:"duration"`
}

func NewProjectHandler(s service.ProjectService) *ProjectHandler {
	return &ProjectHandler{ProjectService: s}
}

// PostProject POST /projects
func (handler *ProjectHandler) PostProject(c echo.Context) error {
	req := PatchProject{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	p := domain.ProjectDetail{
		ID:          id,
		Name:        req.Name,
		Duration:    req.Duration,
		Link:        req.Link,
		Description: req.Description,
		//TODO Members:
	}
	res, err := handler.ProjectService.PostProject(ctx, &p)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, &ProjectResponse{
		ID:       res.ID,
		Name:     res.Name,
		Duration: req.Duration, //TODO
	})
}
