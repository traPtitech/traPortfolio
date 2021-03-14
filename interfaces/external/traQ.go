package external

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type UserResponse struct {
	State       model.TraQState `json:"state"`
	Bot         bool            `json:"bot"`
	DisplayName string          `json:"displayName"`
	Name        string          `json:"name"`
}

type TraQAPI interface {
	GetByID(uuid.UUID, string) (*UserResponse, error)
}
