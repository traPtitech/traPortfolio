//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

type PortalUserResponse struct {
	TraQID         string `json:"id"`
	RealName       string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalAPI interface {
	GetUsers() ([]*PortalUserResponse, error)
	GetUserByTraqID(traQID string) (*PortalUserResponse, error)
}
