package testutils

import (
	"github.com/spf13/viper"
	"github.com/traPtitech/traPortfolio/util/config"
	"testing"
)

func ParseConfig(path string) error {
	viper.AddConfigPath(path)
	return config.ReadFromFile()
}

func GetConfigWithDBName(t *testing.T, dbName string) *config.Config {
	t.Helper()
	return config.Load(func(c *config.Config) {
		c.DB.Name = testDBName(dbName)
	})
}
