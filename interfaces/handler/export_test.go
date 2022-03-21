package handler

import "testing"

func ConvertError(t *testing.T, err error) error {
	t.Helper()

	return convertError(err)
}
