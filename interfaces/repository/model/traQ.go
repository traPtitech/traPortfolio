package model

import "github.com/traPtitech/traPortfolio/domain"

// TraQUser traQ上のユーザー情報
type TraQUser struct {
	State       domain.TraQState
	Bot         bool
	DisplayName string
	Name        string
}
