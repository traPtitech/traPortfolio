//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

type PortalUserResponse struct {
	TraQID         string `json:"id"`
	RealName       string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalAPI interface {
	GetAll() ([]*PortalUserResponse, error)
	GetByTraqID(traQID string) (*PortalUserResponse, error)
}
