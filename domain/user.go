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
	Type        uint
	PrPermitted bool
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
	Name        string
	Result      string
	ContestName string
}

type UserGroup struct {
	ID       uuid.UUID
	Name     string
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
