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

// GetGroups GET /groups
func (h *GroupHandler) GetGroups(_c echo.Context) error {
	c := _c.(*Context)
	ctx := c.Request().Context()
	groups, err := h.srv.GetAllGroups(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Group, len(groups))
	for i, group := range groups {
		res[i] = newGroup(group.ID, group.Name)
	}

	return c.JSON(http.StatusOK, res)
}

// GetGroup GET /groups/:groupID
func (h *GroupHandler) GetGroup(_c echo.Context) error {
	c := _c.(*Context)
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

func formatGetGroup(group *domain.GroupDetail) GroupDetail {
	groupRes := make([]GroupMember, len(group.Members))
	for i, v := range group.Members {
		groupRes[i] = newGroupMember(
			newUser(v.ID, v.Name, v.RealName),
			newYearWithSemesterDuration(
				int(v.Duration.Since.Year),
				int(v.Duration.Since.Semester),
				int(v.Duration.Until.Year),
				int(v.Duration.Until.Semester),
			),
		)
	}
	adminRes := make([]User, len(group.Admin))
	for i, v := range group.Admin {
		adminRes[i] = newUser(v.ID, v.Name, v.RealName)
	}

	res := newGroupDetail(
		newGroup(group.ID, group.Name),
		group.Description,
		adminRes,
		group.Link,
		groupRes,
	)

	return res
}

func newGroupMember(user User, Duration YearWithSemesterDuration) GroupMember {
	return GroupMember{
		Duration: Duration,
		Id:       user.Id,
		Name:     user.Name,
		RealName: user.RealName,
	}
}

func newGroupDetail(group Group, desc string, admin []User, link string, members []GroupMember) GroupDetail {
	return GroupDetail{
		Description: desc,
		Id:          group.Id,
		Admin:       admin,
		Link:        link,
		Members:     members,
		Name:        group.Name,
	}
}
