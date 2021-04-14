package external

import "github.com/gofrs/uuid"

type PortalUserResponse struct {
	TraQID         string `json:"id"`
	RealName       string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalAPI interface {
	GetAll() ([]*PortalUserResponse, error)
	GetByID(id uuid.UUID) (*PortalUserResponse, error)
}
