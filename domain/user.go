package domain

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	RealName string
}

type UserWithDuration struct {
	User     User
	Duration YearWithSemesterDuration
}

type Account struct {
	ID          uuid.UUID
	DisplayName string
	Type        uint8
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
	ID          uuid.UUID // チームID
	Name        string    // チーム名
	Result      string
	ContestName string
}

type UserGroup struct {
	ID       uuid.UUID // Group ID
	Name     string    // Group name
	Duration YearWithSemesterDuration
}

const (
	HOMEPAGE uint8 = iota
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
	TraqStateLimit
)
