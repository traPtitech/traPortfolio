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

type GroupUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Duration domain.GroupDuration
}

// NewGroupUserRespoce creates a GroupUserResponse
func NewGroupUserResponse(id uuid.UUID, name string, dur domain.GroupDuration) *GroupUserResponse {
	return &GroupUserResponse{ID: id, Name: name, Duration: dur}
}

type GroupMemberDetailResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"real_name"`
	Duration domain.GroupDuration
}

// GroupResponse Portfolioのレスポンスで使う班情報
type GroupsResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewGroupResponse creates a GroupHandler
func NewGroupResponse(id uuid.UUID, name string) *GroupsResponse {
	return &GroupsResponse{ID: id, Name: name}
}

func (h *GroupHandler) GetAllGroups(c echo.Context) error {
	ctx := c.Request().Context()
	groups, err := h.srv.GetAllGroups(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]*GroupsResponse, 0, len(groups))
	for _, group := range groups {
		res = append(res, &GroupsResponse{
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
	Leader      *UserResponse
	Members     []*GroupMemberDetailResponse
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
	groupRes := make([]*GroupMemberDetailResponse, 0, len(group.Members))
	for _, v := range group.Members {
		groupRes = append(groupRes, &GroupMemberDetailResponse{
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
		Leader: &UserResponse{
			ID:       group.Leader.ID,
			Name:     group.Leader.Name,
			RealName: group.Leader.RealName,
		},
		Members:     groupRes,
		Description: group.Description,
	}

	return res
}
