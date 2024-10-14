package handler

import (
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"

	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	group repository.GroupRepository
	user  repository.UserRepository
}

// NewGroupHandler creates a GroupHandler
func NewGroupHandler(group repository.GroupRepository, user repository.UserRepository) *GroupHandler {
	return &GroupHandler{group, user}
}

// GetGroups GET /groups
func (h *GroupHandler) GetGroups(c echo.Context) error {
	ctx := c.Request().Context()
	groups, err := h.group.GetGroups(ctx)
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
	group, err := h.group.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// pick all users info
	users, err := h.user.GetUsers(ctx, &repository.GetUsersArgs{}) // TODO: IncludeSuspendedをtrueにするか考える
	if err != nil {
		return err
	}

	umap := make(map[uuid.UUID]*domain.User)
	for _, u := range users {
		umap[u.ID] = u
	}

	// fill members info
	for i, v := range group.Members {
		if u, ok := umap[v.User.ID]; ok {
			m := *u
			group.Members[i].User = m
		}
	}

	// fill leader info
	for i, v := range group.Admin {
		if u, ok := umap[v.ID]; ok {
			m := u
			group.Admin[i] = m
		}
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
		group.Links,
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

func newGroupDetail(group schema.Group, desc string, admin []schema.User, links []string, members []schema.GroupMember) schema.GroupDetail {
	return schema.GroupDetail{
		Description: desc,
		Id:          group.Id,
		Admin:       admin,
		Links:       links,
		Members:     members,
		Name:        group.Name,
	}
}
