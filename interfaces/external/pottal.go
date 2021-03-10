package external

type PortalUserResponse struct {
	ID             string
	Name           string
	AlphabeticName string
}

type PortalAPI interface {
	GetAll() ([]*PortalUserResponse, error)
}
