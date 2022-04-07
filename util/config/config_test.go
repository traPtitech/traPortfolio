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
		DB: config.SQLConfig{
			User: "root",
			Pass: "password",
			Host: "mysql",
			Name: "portfolio",
			Port: 3001,
		},
		Traq: config.TraqConfig{
			Cookie:      "traq cookie",
			APIEndpoint: "traq endpoint",
		},
		Knoq: config.KnoqConfig{
			Cookie:      "knoq cookie",
			APIEndpoint: "knoq endpoint",
		},
		Portal: config.PortalConfig{
			Cookie:      "portal cookie",
			APIEndpoint: "portal endpoint",
		},
	}

	assert.Equal(t, &expected, config.GetConfig())
}
