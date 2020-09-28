package domain

// TraQUser traQ上のユーザー情報
type TraQUser struct {
	State       uint8  `json:"state"` // TODO: 特別な型にする
	Bot         bool   `json:"bot"`
	DisplayName string `json:"displayName"`
	Name        string `json:"name"` // これはportalのIDと一致
}
