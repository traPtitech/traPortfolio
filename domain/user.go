package domain

import (
	"time"

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
	ID        uuid.UUID
	Name      string
	Since     time.Time
	Until     time.Time
	UserSince time.Time
	UserUntil time.Time
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
