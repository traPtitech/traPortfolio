package domain

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	RealName string
}

type Account struct {
	ID          uuid.UUID
	Name        string
	Type        uint
	PrPermitted bool
	URL         string
}

type UserDetail struct {
	User
	State    TraQState
	Bio      string
	Accounts []*Account
}

type UserProject struct {
	ID           uuid.UUID
	Name         string
	Duration     YearWithSemesterDuration
	UserDuration YearWithSemesterDuration
}

type UserContest struct {
	ID          uuid.UUID
	Name        string // チーム名
	Result      string
	ContestName string
}

// UserGroup indicates User who belongs to Group
type UserGroup struct {
	ID       uuid.UUID // User ID
	Name     string    // User Name
	RealName string
	Duration GroupDuration
}

const (
	HOMEPAGE uint = iota
	BLOG
	TWITTER
	FACEBOOK
	PIXIV
	GITHUB
	QIITA
	ATCODER
	SOUNDCLOUD
	AccountLimit
)

type TraQState uint8

const (
	// ユーザーアカウント状態: 凍結
	TraqStateDeactivated TraQState = iota
	// ユーザーアカウント状態: 有効
	TraqStateActive
	// ユーザーアカウント状態: 一時停止
	TraqStateSuspended
)
