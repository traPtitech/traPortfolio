package external

import "github.com/gofrs/uuid"

type UserResponse struct {
	State       uint8  `json:"state"`
	Bot         bool   `json:"bot"`
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
}

type TraqAPI interface {
	GetByID(uuid.UUID, string) (*UserResponse, error)
}
