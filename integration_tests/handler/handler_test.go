package handler

import (
	"testing"

	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
)

func TestMain(m *testing.M) {
	if err := testutils.ParseConfig("../testdata"); err != nil {
		panic(err)
	}

	m.Run()
}
