package handler

import "testing"

// integration_tests/handlerで使えるようにexportしているがあまり綺麗ではない
func ConvertError(t *testing.T, err error) error {
	t.Helper()

	return convertError(err)
}
