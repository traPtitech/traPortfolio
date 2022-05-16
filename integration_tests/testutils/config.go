//go:build integration && db

package testutils

import (
	"github.com/spf13/viper"
	"github.com/traPtitech/traPortfolio/util/config"
)

func ParseConfig(path string) {
	viper.AddConfigPath(path)
	config.ReadFromFile()
}

func GetConfig() *config.Config {
	return config.GetConfig()
}

func GetModified(f config.EditFunc) *config.Config {
	return config.GetModified(f)
}

func GetConfigWithDBName(dbName string) *config.Config {
	return GetModified(func(c *config.Config) {
		c.DB.Name = testDBName(dbName)
	})
}
