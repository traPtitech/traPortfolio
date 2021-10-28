package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	srv service.GroupService
}

type groupParam struct {
	GroupID uuid.UUID `param:"groupID" validate:"is-uuid"`
}

// NewGroupHandler creates a GroupHandler
func NewGroupHandler(service service.GroupService) *GroupHandler {
	return &GroupHandler{service}
}

// GroupResponse Portfolioのレスポンスで使う班情報
type groupsResponse struct {
	ID   uuid.UUID `json:"groupId"`
	Name string    `json:"name"`
}

func (h *GroupHandler) GetAllGroups(c echo.Context) error {
	ctx := c.Request().Context()
	groups, err := h.srv.GetAllGroups(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]*groupsResponse, 0, len(groups))
	for _, group := range groups {
		res = append(res, &groupsResponse{
			ID:   group.ID,
			Name: group.Name,
		})
	}
	return c.JSON(http.StatusOK, res)
}

type groupDetailResponse struct {
	ID          uuid.UUID
	Name        string
	Link        string
	Leader      *domain.User
	Members     []*domain.UserGroup
	Description string
}

func (h *GroupHandler) GetGroup(_c echo.Context) error {
	c := Context{_c}
	req := groupParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	group, err := h.srv.GetGroup(ctx, req.GroupID)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, formatGetGroup(group))
}

func formatGetGroup(group *domain.GroupDetail) *groupDetailResponse {
	groupRes := make([]*domain.UserGroup, 0, len(group.Members))
	for _, v := range group.Members {
		groupRes = append(groupRes, &domain.UserGroup{
			ID:       v.ID,
			Name:     v.Name,
			RealName: v.RealName,
			Duration: domain.GroupDuration{
				Since: domain.YearWithSemester{
					Year:     v.Duration.Since.Year,
					Semester: v.Duration.Since.Semester,
				},
				Until: domain.YearWithSemester{
					Year:     v.Duration.Since.Year,
					Semester: v.Duration.Since.Semester,
				},
			},
		})
	}

	res := &groupDetailResponse{
		ID:   group.ID,
		Name: group.Name,
		Link: group.Link,
		Leader: &domain.User{
			ID:       group.Leader.ID,
			Name:     group.Leader.Name,
			RealName: group.Leader.RealName,
		},
		Members:     groupRes,
		Description: group.Description,
	}

	return res
}
