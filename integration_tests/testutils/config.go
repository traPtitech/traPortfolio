package testutils

import (
	"testing"

	"github.com/traPtitech/traPortfolio/util/config"
)

func GetConfigWithDBName(t *testing.T, dbName string) *config.Config {
	t.Helper()
	return config.Load(func(c *config.Config) {
		c.DB.Name = testDBName(dbName)
	})
}
