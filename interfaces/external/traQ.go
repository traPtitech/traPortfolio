package external

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type TraQUserResponse struct {
	State       domain.TraQState `json:"state"`
	Bot         bool             `json:"bot"`
	DisplayName string           `json:"displayName"`
	Name        string           `json:"name"`
}

type GroupUserResponse struct {
	ID       uuid.UUID `json:"groupId"`
	Name     string    `json:"name"`
	Duration domain.ProjectDuration
}

type GroupsResponse struct {
	ID   uuid.UUID `json:"Id"`
	Name string    `json:"Name"`
}

type GroupDetailResponse struct {
	ID          uuid.UUID `json:"groupId"`
	Name        string    `json:"name"`
	Link        string    `json:"link"`
	Leader      domain.User
	Members     []domain.UserGroup
	Description string `json:"description"`
}

type TraQAPI interface {
	GetByID(id uuid.UUID) (*TraQUserResponse, error)
}

type GroupAPI interface {
	GetAllGroups() ([]*GroupsResponse, error)
	GetGroup(groupID uuid.UUID) (*GroupDetailResponse, error)
}
