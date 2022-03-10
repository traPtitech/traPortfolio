//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type TraQUserResponse struct {
	ID    uuid.UUID        `json:"id"`
	State domain.TraQState `json:"state"`
}

type TraQGetAllArgs struct {
	IncludeSuspended bool
	Name             string
}

type TraQAPI interface {
	GetAll(args *TraQGetAllArgs) ([]*TraQUserResponse, error)
	GetByID(id uuid.UUID) (*TraQUserResponse, error)
}
