package random

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
	"unsafe"

	"github.com/gofrs/uuid"
)

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

// AlphaNumeric 指定した文字数のランダム英数字文字列を生成します
// この関数はmath/randが生成する擬似乱数を使用します
func AlphaNumeric(n int) string {
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

// UUID ランダムなUUIDを生成します
func UUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func Time() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).UnixMicro()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).UnixMicro()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.UnixMicro(sec).In(time.UTC)
}

func URL(useHTTPS bool, domainLength uint16) *url.URL {
	scheme := "https"
	if !useHTTPS {
		scheme = "http"
	}
	scheme += "://"

	scheme += fmt.Sprintf("%s%s", scheme, AlphaNumeric(int(domainLength)))
	url, err := url.Parse(scheme)
	if err != nil {
		panic(err)
	}
	return url
}

func RandURLString() string {
	return URL(rand.Intn(2) < 1, uint16(rand.Intn(20)+1)).String()
}
