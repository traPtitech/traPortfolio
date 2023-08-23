package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	s service.GroupService
}

// NewGroupHandler creates a GroupHandler
func NewGroupHandler(s service.GroupService) *GroupHandler {
	return &GroupHandler{s}
}

// GetGroups GET /groups
func (h *GroupHandler) GetGroups(c echo.Context) error {
	ctx := c.Request().Context()
	groups, err := h.s.GetAllGroups(ctx)
	if err != nil {
		return err
	}

	res := make([]schema.Group, len(groups))
	for i, group := range groups {
		res[i] = newGroup(group.ID, group.Name)
	}

	return c.JSON(http.StatusOK, res)
}

// GetGroup GET /groups/:groupID
func (h *GroupHandler) GetGroup(c echo.Context) error {
	groupID, err := getID(c, keyGroupID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	group, err := h.s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, formatGetGroup(group))
}

func formatGetGroup(group *domain.GroupDetail) schema.GroupDetail {
	groupRes := make([]schema.GroupMember, len(group.Members))
	for i, v := range group.Members {
		groupRes[i] = newGroupMember(
			newUser(v.User.ID, v.User.Name, v.User.RealName()),
			schema.ConvertDuration(v.Duration),
		)
	}
	adminRes := make([]schema.User, len(group.Admin))
	for i, v := range group.Admin {
		adminRes[i] = newUser(v.ID, v.Name, v.RealName())
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

func newGroupMember(user schema.User, Duration schema.YearWithSemesterDuration) schema.GroupMember {
	return schema.GroupMember{
		Duration: Duration,
		Id:       user.Id,
		Name:     user.Name,
		RealName: user.RealName,
	}
}

func newGroupDetail(group schema.Group, desc string, admin []schema.User, link string, members []schema.GroupMember) schema.GroupDetail {
	return schema.GroupDetail{
		Description: desc,
		Id:          group.Id,
		Admin:       admin,
		Link:        link,
		Members:     members,
		Name:        group.Name,
	}
}
