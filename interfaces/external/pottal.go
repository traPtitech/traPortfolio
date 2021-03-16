package external

type PortalUserResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalAPI interface {
	GetAll() ([]*PortalUserResponse, error)
}
