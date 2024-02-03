package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	viper.AddConfigPath("./testdata")
	assert.NoError(t, Parse())

	expected := Config{
		IsProduction:   true,
		Port:           3000,
		OnlyMigrate:    true,
		InsertMockData: true,
		DB: SQLConfig{
			User:    "root",
			Pass:    "password",
			Host:    "mysql",
			Name:    "portfolio",
			Port:    3001,
			Verbose: true,
		},
		Traq: APIConfig{
			Cookie:      "traq cookie",
			APIEndpoint: "traq endpoint",
		},
		Knoq: APIConfig{
			Cookie:      "knoq cookie",
			APIEndpoint: "knoq endpoint",
		},
		Portal: APIConfig{
			Cookie:      "portal cookie",
			APIEndpoint: "portal endpoint",
		},
	}

	assert.Equal(t, &expected, Load())
}
