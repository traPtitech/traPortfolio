package domain

import (
	"testing"
)

func (u User) RealNameForTest(t *testing.T) string {
	t.Helper()
	return u.realName
}
