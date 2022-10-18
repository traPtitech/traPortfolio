package random

import (
	"math/rand"
	"net/url"
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

func Duration() domain.YearWithSemesterDuration {
	ys1 := domain.YearWithSemester{
		Year:     Time().Year(),
		Semester: rand.Intn(2),
	}
	ys2 := domain.YearWithSemester{
		Year:     Time().Year(),
		Semester: rand.Intn(2),
	}

	if ys1.After(ys2) {
		return domain.YearWithSemesterDuration{
			Since: ys2,
			Until: ys1,
		}
	}

	return domain.YearWithSemesterDuration{
		Since: ys1,
		Until: ys2,
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
