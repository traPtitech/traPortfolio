package testutils

import (
	"github.com/spf13/viper"
	"github.com/traPtitech/traPortfolio/util/config"
)

func ParseConfig(path string) {
	viper.AddConfigPath(path)
	config.Parse()
}

func GetConfig() *config.Config {
	return config.GetConfig()
}
