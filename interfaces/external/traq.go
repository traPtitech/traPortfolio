//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

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
	GetByID(id uuid.UUID) (*TraQUserResponse, error)
}
