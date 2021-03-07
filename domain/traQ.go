package domain

type TraQState uint8

// TraQUser traQ上のユーザー情報
type TraQUser struct {
	State       TraQState
	Bot         bool
	DisplayName string
	Name        string
}
