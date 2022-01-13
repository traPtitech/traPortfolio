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

// type GroupUserResponse struct {
// 	ID       uuid.UUID `json:"id"`
// 	Name     string    `json:"name"`
// 	Duration domain.GroupDuration
// }

// NewGroupUserRespoce creates a GroupUserResponse
// func NewGroupUserResponse(id uuid.UUID, name string, dur domain.GroupDuration) *GroupUserResponse {
// 	return &GroupUserResponse{ID: id, Name: name, Duration: dur}
// }

// type GroupMemberDetailResponse struct {
// 	ID       uuid.UUID `json:"id"`
// 	Name     string    `json:"name"`
// 	RealName string    `json:"real_name"`
// 	Duration domain.GroupDuration
// }

// // GroupResponse Portfolioのレスポンスで使う班情報
// type GroupsResponse struct {
// 	ID   uuid.UUID `json:"id"`
// 	Name string    `json:"name"`
// }

// NewGroupResponse creates a GroupHandler
// func NewGroupResponse(id uuid.UUID, name string) *GroupsResponse {
// 	return &GroupsResponse{ID: id, Name: name}
// }

func (h *GroupHandler) GetAllGroups(c echo.Context) error {
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

// type groupDetailResponse struct {
// 	ID          uuid.UUID
// 	Name        string
// 	Link        string
// 	Leader      *UserResponse
// 	Members     []*GroupMemberDetailResponse
// 	Description string
// }

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
