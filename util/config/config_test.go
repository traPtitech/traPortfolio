package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
		Traq: TraqConfig{
			AccessToken: "traq token",
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

	viper.AddConfigPath("./testdata")

	got, err := Load(LoadOpts{})
	assert.NoError(t, err)
	assert.Equal(t, &expected, got)
}
