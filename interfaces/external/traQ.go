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

type TraQAPI interface {
	GetByID(uuid.UUID) (*TraQUserResponse, error)
}
