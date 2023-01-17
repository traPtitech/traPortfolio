package random

import (
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"time"
	"unsafe"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

// AlphaNumericn 指定した文字数のランダム英数字文字列を生成します
// この関数はmath/randが生成する擬似乱数を使用します
func AlphaNumericn(n int) string {
	b := make([]byte, n)
	cache, remain := rand.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func AlphaNumeric() string {
	return AlphaNumericn(rand.Intn(30) + 1)
}

// UUID ランダムなUUIDを生成します
func UUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func SinceAndUntil() (time.Time, time.Time) {
	since := Time()
	until := Time()

	if since.After(until) {
		since, until = until, since
	}

	return since, until
}

func Time() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).UnixMicro()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).UnixMicro()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.UnixMicro(sec).In(time.UTC)
}

func URL(useHTTPS bool, domainLength int) *url.URL {
	scheme := "https"
	if !useHTTPS {
		scheme = "http"
	}
	scheme += "://"

	scheme += AlphaNumericn(domainLength)
	url, err := url.Parse(scheme)
	if err != nil {
		panic(err)
	}
	return url
}

func RandURLString() string {
	return URL(rand.Intn(2) < 1, rand.Intn(20)+1).String()
}

func AccountURLString(accountType uint) string {
	var AccountURLs = map[uint][]string{
		domain.TWITTER: {
			"https://twitter.com/qi_1WI_nku",
			"https://twitter.com/XF1G6_kqEG",
			"https://twitter.com/Aer7qyNEUz",
		},
		domain.FACEBOOK: {
			"https://www.facebook.com/hYFmBZH21e",
			"https://www.facebook.com/3LwQ.v7uyN",
			"https://www.facebook.com/c6xegdX.hu",
		},
		domain.PIXIV: {
			"https://www.pixiv.net/users/2102945291",
			"https://www.pixiv.net/users/4818932326",
			"https://www.pixiv.net/users/1939586271",
		},
		domain.GITHUB: {
			"https://github.com/WeL-rKj-xz",
			"https://github.com/q7DO-T9GTO",
			"https://github.com/wexjxusr1B",
		},
		domain.QIITA: {
			"https://qiita.com/BSpYu_LyYg",
			"https://qiita.com/HV6-Ik252Z",
			"https://qiita.com/5XcnQ8fyze",
		},
		domain.ZENN: {
			"https://zenn.dev/2Kl1M.I3MO",
			"https://zenn.dev/we.Xh9Sg2k",
			"https://zenn.dev/ygZsTx1Pjf",
		},
		domain.ATCODER: {
			"https://atcoder.jp/users/Ib_ucf2TjO",
			"https://atcoder.jp/users/8d_z3Dm_T1",
			"https://atcoder.jp/users/yKUfEWAnNB",
		},
		domain.SOUNDCLOUD: {
			"https://soundcloud.com/ofb4igxvi8",
			"https://soundcloud.com/r_e-dt6qgn",
			"https://soundcloud.com/zut7-ajedl",
		},
		domain.HACKTHEBOX: {
			"https://app.hackthebox.com/users/ORIuZ5qoXl",
			"https://app.hackthebox.com/users/8WIScK32pB",
			"https://app.hackthebox.com/users/IuqV2A1ux1",
		},
		domain.CTFTIME: {
			"https://ctftime.org/user/1939138413",
			"https://ctftime.org/user/4285429253",
			"https://ctftime.org/user/8295210365",
		},
	}
	if accountType == domain.HOMEPAGE || accountType == domain.BLOG {
		return fmt.Sprintf("https://%s", AlphaNumeric())
	}
	return AccountURLs[accountType][rand.Intn(3)]
}

func Duration() domain.YearWithSemesterDuration {
	yss := []domain.YearWithSemester{
		{
			Year:     Time().Year(),
			Semester: rand.Intn(2),
		},
		{
			Year:     Time().Year(),
			Semester: rand.Intn(2),
		},
	}

	// 時系列昇順に並べる
	sort.Slice(yss, func(i, j int) bool {
		return !yss[i].After(yss[j])
	})

	return domain.YearWithSemesterDuration{
		Since: yss[0],
		Until: yss[1],
	}
}

func Uint8n(n uint8) uint8 {
	return uint8(rand.Int31n(int32(n)))
}

func Bool() bool {
	return rand.Int()%2 == 0
}

func OptBool() optional.Bool {
	return optional.NewBool(Bool(), Bool())
}

func OptBoolNotNull() optional.Bool {
	return optional.NewBool(Bool(), true)
}

func OptInt64() optional.Int64 {
	return optional.NewInt64(rand.Int63(), Bool())
}

func OptInt64n(n int64) optional.Int64 {
	return optional.NewInt64(rand.Int63n(n), Bool())
}

func OptInt64nNotNull(n int64) optional.Int64 {
	return optional.NewInt64(rand.Int63n(n), true)
}

func OptAlphaNumericn(n int) optional.String {
	return optional.NewString(AlphaNumericn(n), Bool())
}

func OptAlphaNumeric() optional.String {
	return optional.NewString(AlphaNumeric(), Bool())
}

func OptAlphaNumericNotNull() optional.String {
	return optional.NewString(AlphaNumeric(), true)
}

func OptTime() optional.Time {
	return optional.NewTime(Time(), Bool())
}

func OptURLString() optional.String {
	return optional.NewString(RandURLString(), Bool())
}

func OptURLStringNotNull() optional.String {
	return optional.NewString(RandURLString(), true)
}

func OptAccountURLStringNotNull(accountType uint) optional.String {
	return optional.NewString(AccountURLString(accountType), true)
}
