package external

type PortalUserResponse struct {
	TraQID         string `json:"id"`
	RealName       string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalAPI interface {
	GetAll() ([]*PortalUserResponse, error)
	GetByID(string) (*PortalUserResponse, error)
}
