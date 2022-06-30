package handler

import (
	"testing"

	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
)

func TestMain(m *testing.M) {
	testutils.ParseConfig("../testdata")

	m.Run()
}
