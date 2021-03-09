package external

type PortalUserResponse struct {
	ID             string
	Name           string
	AlphabeticName string
}

type PortalAPI interface {
	GetAll(string) ([]*PortalUserResponse, error)
	GetByID(string, string) (*PortalUserResponse, error)
}
