package domain

// PortalUser Portal上のユーザー情報
type PortalUser struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}
