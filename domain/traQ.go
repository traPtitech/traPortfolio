package domain

type TraQState uint8

const (
	// ユーザーアカウント状態: 凍結
	TraqStateDeactivated TraQState = iota
	// ユーザーアカウント状態: 有効
	TraqStateActive
	// ユーザーアカウント状態: 一時停止
	TraqStateSuspended
	TraqStateLimit
)

// TraQUser traQ上のユーザー情報
type TraQUser struct {
	State       TraQState
	Bot         bool
	DisplayName string
	Name        string
}
