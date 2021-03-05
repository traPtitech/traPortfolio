package domain

type State uint8

// TraQUser traQ上のユーザー情報
type TraQUser struct {
	State
	Bot         bool
	DisplayName string
	Name        string
}
