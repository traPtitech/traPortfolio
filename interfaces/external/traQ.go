package external

import "github.com/gofrs/uuid"

type UserResponse struct {
	State       uint8
	Bot         bool
	DisplayName string
	Name        string
}

type TraqAPI interface {
	GetByID(uuid.UUID, string) (*UserResponse, error)
}
