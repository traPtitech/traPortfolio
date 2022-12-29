package domain

import (
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	RealName string // TODO: 後で消す
	Check    bool
}

func NewUser(id uuid.UUID, name string, realName string, check bool) *User {
	return &User{
		ID:       id,
		Name:     name,
		RealName: realName,
		Check:    check,
	}
}

type UserWithDuration struct {
	User     User
	Duration YearWithSemesterDuration
}

type Account struct {
	ID          uuid.UUID
	DisplayName string
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
	ID        uuid.UUID // コンテストID
	Name      string    // コンテスト名
	TimeStart time.Time
	TimeEnd   time.Time
	Teams     []*ContestTeam // ユーザーが所属するチームのリスト
}

type UserGroup struct {
	ID       uuid.UUID // Group ID
	Name     string    // Group name
	Duration YearWithSemesterDuration
}

type AccountURL struct {
	URL    string
	Regexp string
}

type AccountType uint8

func IsValidAccountURL(accountType AccountType, URL string) bool {
	if r, ok := urlRegexp[uint(accountType)]; ok {
		return r.MatchString(URL)
	}
	return false
}

const (
	HOMEPAGE uint = iota
	BLOG
	TWITTER
	FACEBOOK
	PIXIV
	GITHUB
	QIITA
	ZENN
	ATCODER
	SOUNDCLOUD
	HACKTHEBOX
	CTFTIME
	AccountLimit
)

var urlRegexp = map[uint]*regexp.Regexp{
	HOMEPAGE:   regexp.MustCompile("^https?://.+$"),
	BLOG:       regexp.MustCompile("^https?://.+$"),
	TWITTER:    regexp.MustCompile("^https?://twitter.com/[a-zA-Z0-9_]+$"),
	FACEBOOK:   regexp.MustCompile("^https?://www.facebook.com/[a-zA-Z0-9.]+$"),
	PIXIV:      regexp.MustCompile("^https?://www.pixiv.net/users/[0-9]+"),
	GITHUB:     regexp.MustCompile("^https?://github.com/[a-zA-Z0-9-]+$"),
	QIITA:      regexp.MustCompile("^https?://qiita.com/[a-zA-Z0-9-_]+$"),
	ZENN:       regexp.MustCompile("^https?://zenn.dev/[a-zA-Z0-9.]+$"),
	ATCODER:    regexp.MustCompile("^https?://atcoder.jp/users/[a-zA-Z0-9_]+$"),
	SOUNDCLOUD: regexp.MustCompile("^https?://soundcloud.com/[a-z0-9-_]+$"),
	HACKTHEBOX: regexp.MustCompile("^https?://app.hackthebox.com/users/[a-zA-Z0-9]+$"),
	CTFTIME:    regexp.MustCompile("^https?://ctftime.org/user/[0-9]+$"),
}

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
