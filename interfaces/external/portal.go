package external

import "github.com/gofrs/uuid"

type PortalUserResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalQAPI interface {
	GetAll() ([]*PortalUserResponse, error)
	GetByID(id uuid.UUID) (PortalUserResponse, error)
}
