package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/util/config"
)

func TestLoad(t *testing.T) {
	defaultConfig := config.Config{
		IsProduction:   false,
		Port:           1323,
		OnlyMigrate:    false,
		InsertMockData: false,
		DB: config.SQLConfig{
			User:    "root",
			Pass:    "password",
			Host:    "localhost",
			Name:    "portfolio",
			Port:    3306,
			Verbose: false,
		},
		Traq: config.TraqConfig{
			AccessToken: "",
		},
		Knoq: config.APIConfig{
			Cookie:      "",
			APIEndpoint: "",
		},
		Portal: config.APIConfig{
			Cookie:      "",
			APIEndpoint: "",
		},
	}

	t.Run("default", func(t *testing.T) {
		got, err := config.Load(config.LoadOpts{})
		assert.NoError(t, err)
		assert.Equal(t, &defaultConfig, got)
	})

	t.Run("from file", func(t *testing.T) {
		yaml := `
production: true
port: 8000`
		configPath := filepath.Join(t.TempDir(), "config.yaml")
		os.Create(configPath)
		os.WriteFile(configPath, []byte(yaml), 0644)
		t.Setenv("TPF_CONFIG", configPath)

		expected := defaultConfig
		expected.IsProduction = true
		expected.Port = 8000

		got, err := config.Load(config.LoadOpts{})
		assert.NoError(t, err)
		assert.Equal(t, &expected, got)
	})

	t.Run("from env", func(t *testing.T) {
		t.Setenv("TPF_PRODUCTION", "true")
		t.Setenv("TPF_PORT", "8000")

		expected := defaultConfig
		expected.IsProduction = true
		expected.Port = 8000

		got, err := config.Load(config.LoadOpts{})
		assert.NoError(t, err)
		assert.Equal(t, &expected, got)
	})

	t.Run("from flag", func(t *testing.T) {
		pflag.CommandLine.Set("production", "true")
		pflag.CommandLine.Set("port", "8000")
		t.Cleanup(func() {
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
		})

		expected := defaultConfig
		expected.IsProduction = true
		expected.Port = 8000

		got, err := config.Load(config.LoadOpts{})
		assert.NoError(t, err)
		assert.Equal(t, &expected, got)
	})

	t.Run("priority order is flag, env, file, then default", func(t *testing.T) {
		t.Skip("It fails if flag is set twice")

		yaml := `
db:
  user: file
  pass: file
  host: file`
		configPath := filepath.Join(t.TempDir(), "config.yaml")
		os.Create(configPath)
		os.WriteFile(configPath, []byte(yaml), 0644)
		t.Setenv("TPF_CONFIG", configPath)

		t.Setenv("TPF_DB_USER", "env")
		t.Setenv("TPF_DB_PASS", "env")

		pflag.CommandLine.Set("db-user", "flag")
		t.Cleanup(func() {
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
		})

		expected := defaultConfig
		expected.DB.User = "flag"
		expected.DB.Pass = "env"
		expected.DB.Host = "file"

		got, err := config.Load(config.LoadOpts{})
		assert.NoError(t, err)
		assert.Equal(t, &expected, got)
	})
}
