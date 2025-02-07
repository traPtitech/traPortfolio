package domain

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	realName string
	Check    bool
}

func NewUser(id uuid.UUID, name string, realName string, check bool) *User {
	return &User{
		ID:       id,
		Name:     name,
		realName: realName,
		Check:    check,
	}
}

func (u User) RealName() string {
	if !u.Check {
		return ""
	}

	return u.realName
}

type UserWithDuration struct {
	User     User
	Duration YearWithSemesterDuration
}

type Account struct {
	ID          uuid.UUID
	DisplayName string
	Type        AccountType
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
	Teams     []*ContestTeamWithoutMembers // ユーザーが所属するチームのリスト
}

type UserGroup struct {
	ID       uuid.UUID // Group ID
	Name     string    // Group name
	Duration YearWithSemesterDuration
}

type AccountType uint8

var (
	_ sql.Scanner   = (*AccountType)(nil)
	_ driver.Valuer = AccountType(0)
)

const (
	HOMEPAGE AccountType = iota
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
	BLUESKY
	MIXI2
	AccountLimit
)

func (a *AccountType) Scan(src interface{}) error {
	s := sql.NullByte{}
	if err := s.Scan(src); err != nil {
		return err
	}

	if s.Valid {
		newAT := AccountType(s.Byte)
		if newAT >= AccountLimit {
			return fmt.Errorf("%w: AccountType(%d) must be less than %d", ErrTooLargeEnum, newAT, AccountLimit)
		}

		*a = newAT
	}

	return nil
}

func (a AccountType) Value() (driver.Value, error) {
	return sql.NullByte{Byte: byte(a), Valid: true}.Value()
}

func IsValidAccountURL(accountType AccountType, URL string) bool {
	if _, err := url.Parse(URL); err != nil {
		return false
	}

	var urlRegexp = map[AccountType]*regexp.Regexp{
		HOMEPAGE:   regexp.MustCompile(`^https?://.+$`),
		BLOG:       regexp.MustCompile(`^https?://.+$`),
		TWITTER:    regexp.MustCompile(`^https://(twitter|x)\.com/[a-zA-Z0-9_]+$`),
		FACEBOOK:   regexp.MustCompile(`^https://www\.facebook\.com/[a-zA-Z0-9.]+$`),
		PIXIV:      regexp.MustCompile(`^https://www\.pixiv\.net/users/[0-9]+`),
		GITHUB:     regexp.MustCompile(`^https://github\.com/[a-zA-Z0-9-]+$`),
		QIITA:      regexp.MustCompile(`^https://qiita\.com/[a-zA-Z0-9-_]+$`),
		ZENN:       regexp.MustCompile(`^https://zenn\.dev/[a-z0-9_]+$`),
		ATCODER:    regexp.MustCompile(`^https://atcoder\.jp/users/[a-zA-Z0-9_]+$`),
		SOUNDCLOUD: regexp.MustCompile(`^https://soundcloud\.com/[a-z0-9-_]+$`),
		HACKTHEBOX: regexp.MustCompile(`^https://app\.hackthebox\.com/users/[a-zA-Z0-9]+$`),
		CTFTIME:    regexp.MustCompile(`^https://ctftime\.org/user/[0-9]+$`),
		BLUESKY:    regexp.MustCompile(`^https://bsky\.app/profile/[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](\.[a-zA-Z0-9]+)+$`),
		MIXI2:      regexp.MustCompile(`^https://mixi\.social/@[a-zA-Z][a-zA-Z0-9_]{2,14}[a-zA-Z0-9]$`),
	}

	if r, ok := urlRegexp[accountType]; ok {
		return r.MatchString(URL)
	}

	return false
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
