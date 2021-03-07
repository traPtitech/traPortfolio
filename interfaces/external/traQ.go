package external

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type UserResponse struct {
	State       domain.TraQState `json:"state"`
	Bot         bool             `json:"bot"`
	DisplayName string           `json:"displayName"`
	Name        string           `json:"name"`
}

type TraqAPI interface {
	GetByID(uuid.UUID, string) (*UserResponse, error)
}
