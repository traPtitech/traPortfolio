package domain

import (
	"fmt"
	"regexp"
	"time"

	vd "github.com/go-ozzo/ozzo-validation/v4"
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

func (t AccountType) URLValidate() vd.MatchRule {
	regexpText := fmt.Sprintf("^https?://%s/%s$", URLRegexp[uint(t)].URL, URLRegexp[uint(t)].Regexp)
	if t == AccountType(HOMEPAGE) || t == AccountType(BLOG) {
		return vd.Match(regexp.MustCompile("^https?://.+$"))
	} else {
		return vd.Match(regexp.MustCompile(regexpText))
	}
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

var URLRegexp = map[uint]AccountURL{
	HOMEPAGE:   {"", ""},
	BLOG:       {"", ""},
	TWITTER:    {"twitter.com", "[a-zA-Z0-9_]+"},
	FACEBOOK:   {"www.facebook.com", "[a-zA-Z0-9.]+"},
	PIXIV:      {"www.pixiv.net/users", "[0-9]+"},
	GITHUB:     {"github.com", "[a-zA-Z0-9-]+"},
	QIITA:      {"qiita.com", "[a-zA-Z0-9-_]+"},
	ZENN:       {"zenn.dev", "[a-zA-Z0-9.]+"},
	ATCODER:    {"atcoder.jp/users", "[a-zA-Z0-9_]+"},
	SOUNDCLOUD: {"soundcloud.com", "[a-z0-9-_]+"},
	HACKTHEBOX: {"app.hackthebox.com/users", "[a-zA-Z0-9]+"},
	CTFTIME:    {"ctftime.org/user", "[0-9]+"},
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
