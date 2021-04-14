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
	Type        uint
	PrPermitted bool
}

type UserDetail struct {
	ID       uuid.UUID
	Name     string
	RealName string
	State    TraQState
	Bio      string
	Accounts []*Account
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
