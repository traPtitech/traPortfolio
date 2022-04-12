package config_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/util/config"
)

func TestParse(t *testing.T) {
	viper.AddConfigPath("./testdata")
	config.Parse()

	expected := config.Config{
		IsProduction: true,
		Port:         3000,
		Migrate:      true,
		DB: struct {
			User string `mapstructure:"user"`
			Pass string `mapstructure:"pass"`
			Host string `mapstructure:"host"`
			Name string `mapstructure:"name"`
			Port int    `mapstructure:"port"`
		}{
			User: "root",
			Pass: "password",
			Host: "mysql",
			Name: "portfolio",
			Port: 3001,
		},
		Traq: struct {
			Cookie      string `mapstructure:"cookie"`
			APIEndpoint string `mapstructure:"apiEndpoint"`
		}{
			Cookie:      "traq cookie",
			APIEndpoint: "traq endpoint",
		},
		Knoq: struct {
			Cookie      string `mapstructure:"cookie"`
			APIEndpoint string `mapstructure:"apiEndpoint"`
		}{
			Cookie:      "knoq cookie",
			APIEndpoint: "knoq endpoint",
		},
		Portal: struct {
			Cookie      string `mapstructure:"cookie"`
			APIEndpoint string `mapstructure:"apiEndpoint"`
		}{
			Cookie:      "portal cookie",
			APIEndpoint: "portal endpoint",
		},
	}

	assert.Equal(t, &expected, config.GetConfig())
}
