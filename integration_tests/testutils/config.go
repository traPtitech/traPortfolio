package testutils

import (
	"github.com/spf13/viper"
	"github.com/traPtitech/traPortfolio/util/config"
)

func ParseConfig(path string) error {
	viper.AddConfigPath(path)
	return config.ReadFromFile()
}

func GetConfigWithDBName(dbName string) *config.Config {
	return config.Load(func(c *config.Config) {
		c.DB.Name = testDBName(dbName)
	})
}
