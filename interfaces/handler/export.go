//go:build integration && db

// NOTE: パッケージを跨ぐためbuild tagsを使ってexportしている
package handler

import "testing"

func ConvertError(t *testing.T, err error) error {
	t.Helper()

	return convertError(err)
}
