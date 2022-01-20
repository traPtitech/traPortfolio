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

func (h *GroupHandler) GetAllGroups(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	groups, err := h.srv.GetAllGroups(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]Group, 0, len(groups))
	for i, group := range groups {
		res[i] = newGroup(group.ID, group.Name)
	}

	return c.JSON(http.StatusOK, res)
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

func formatGetGroup(group *domain.GroupDetail) GroupDetail {
	groupRes := make([]GroupMember, 0, len(group.Members))
	for i, v := range group.Members {
		groupRes[i] = newGroupMember(
			newUser(v.ID, v.Name, v.RealName),
			newYearWithSemesterDuration(
				YearWithSemester{
					Semester: Semester(v.Duration.Since.Semester),
					Year:     int(v.Duration.Since.Year),
				},
				YearWithSemester{
					Semester: Semester(v.Duration.Until.Semester),
					Year:     int(v.Duration.Until.Year),
				},
			),
		)
	}

	res := newGroupDetail(
		newGroup(group.ID, group.Name),
		group.Description,
		newUser(group.Leader.ID, group.Leader.Name, group.Leader.RealName),
		group.Link,
		groupRes,
	)

	return res
}

func newGroupMember(user User, Duration YearWithSemesterDuration) GroupMember {
	return GroupMember{
		User:     user,
		Duration: Duration,
	}
}

func newGroupDetail(group Group, desc string, leader User, link string, members []GroupMember) GroupDetail {
	return GroupDetail{
		Group:       group,
		Description: desc,
		Leader:      leader,
		Link:        link,
		Members:     members,
	}
}
