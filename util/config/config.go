package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	config   Config
	isParsed atomic.Bool

	flagKeys = []struct{ path, flag string }{
		{"production", "production"},
		{"port", "port"},
		{"onlyMigrate", "only-migrate"},
		{"insertMockData", "insert-mock-data"},
		{"db.user", "db-user"},
		{"db.pass", "db-pass"},
		{"db.host", "db-host"},
		{"db.name", "db-name"},
		{"db.port", "db-port"},
		{"db.verbose", "db-verbose"},
		{"traq.cookie", "traq-cookie"},
		{"traq.apiEndpoint", "traq-api-endpoint"},
		{"knoq.cookie", "knoq-cookie"},
		{"knoq.apiEndpoint", "knoq-api-endpoint"},
		{"portal.cookie", "portal-cookie"},
		{"portal.apiEndpoint", "portal-api-endpoint"},
	}
)

type (
	// Immutable except within this package or EditFunc
	Config struct {
		IsProduction   bool `mapstructure:"production"`
		Port           int  `mapstructure:"port"`
		OnlyMigrate    bool `mapstructure:"onlyMigrate"`
		InsertMockData bool `mapstructure:"insertMockData"`

		DB     SQLConfig    `mapstructure:"db"`
		Traq   TraqConfig   `mapstructure:"traq"`
		Knoq   KnoqConfig   `mapstructure:"knoq"`
		Portal PortalConfig `mapstructure:"portal"`
	}

	SQLConfig struct {
		User    string `mapstructure:"user"`
		Pass    string `mapstructure:"pass"`
		Host    string `mapstructure:"host"`
		Name    string `mapstructure:"name"`
		Port    int    `mapstructure:"port"`
		Verbose bool   `mapstructure:"verbose"`
	}

	// NOTE: wireが複数の同じ型の変数を扱えないためdefined typeを用いる
	// Ref: https://github.com/google/wire/blob/d07cde0df9/docs/faq.md#what-if-my-dependency-graph-has-two-dependencies-of-the-same-type
	TraqConfig   APIConfig
	KnoqConfig   APIConfig
	PortalConfig APIConfig

	APIConfig struct {
		Cookie      string `mapstructure:"cookie"`
		APIEndpoint string `mapstructure:"apiEndpoint"`
	}

	EditFunc func(*Config)
)

func init() {
	isParsed.Store(false)

	pflag.Bool("production", false, "whether production or development")
	pflag.Int("port", 1323, "api port")
	pflag.Bool("only-migrate", false, "only migrate db (not start server)")
	pflag.Bool("insert-mock-data", false, "insert sample mock data(for dev)")

	pflag.String("db-user", "root", "db user name")
	pflag.String("db-pass", "password", "db password")
	pflag.String("db-host", "localhost", "db host")
	pflag.String("db-name", "portfolio", "db name")
	pflag.Int("db-port", 3306, "db port")
	pflag.Bool("db-verbose", false, "db verbose mode")
	pflag.String("traq-cookie", "", "traq cookie")
	pflag.String("traq-api-endpoint", "", "traq api endpoint")
	pflag.String("knoq-cookie", "", "knoq cookie")
	pflag.String("knoq-api-endpoint", "", "knoq api endpoint")
	pflag.String("portal-cookie", "", "portal cookie")
	pflag.String("portal-api-endpoint", "", "portal api endpoint")
	pflag.StringP("config", "c", "", "config file path")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func Parse() error {
	pflag.Parse()
	if err := ReadFromFile(); err != nil {
		return fmt.Errorf("read config from file: %w", err)
	}

	return nil
}

func ReadFromFile() error {
	for _, key := range flagKeys {
		if err := viper.BindPFlag(key.path, pflag.Lookup(key.flag)); err != nil {
			return fmt.Errorf("bind flag %s: %w", key.flag, err)
		}
	}

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("bind flags: %w", err)
	}

	configPath, err := pflag.CommandLine.GetString("config")
	if err != nil {
		return fmt.Errorf("get config flag: %w", err)
	}

	if len(configPath) > 0 {
		viper.SetConfigFile(configPath)
	} else {
		// default path is ./config.yaml
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if len(configPath) > 0 {
				return fmt.Errorf("read config from %s: %w", configPath, err)
			}

			log.Printf("config file does not found: %v", err)
		} else {
			return fmt.Errorf("read config: %w", err)
		}
	} else {
		log.Printf("successfully loaded config from %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	isParsed.Store(true)

	return nil
}

func ReadDefault() error {
	for _, key := range flagKeys {
		if err := viper.BindPFlag(key.path, pflag.Lookup(key.flag)); err != nil {
			return fmt.Errorf("bind flag %s: %w", key.flag, err)
		}
	}

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("bind flags: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	isParsed.Store(true)

	return nil
}

func Load(editFuncs ...EditFunc) *Config {
	if !isParsed.Load() {
		panic("config does not parsed")
	}

	cloned := config.clone()
	for _, f := range editFuncs {
		f(cloned)
	}

	return cloned
}

func (c *Config) clone() *Config {
	cloned := *c
	return &cloned
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *Config) TraqConf() *TraqConfig {
	return &c.Traq
}

func (c *Config) KnoqConf() *KnoqConfig {
	return &c.Knoq
}

func (c *Config) PortalConf() *PortalConfig {
	return &c.Portal
}

func (c *TraqConfig) API() *APIConfig {
	return (*APIConfig)(c)
}

func (c *KnoqConfig) API() *APIConfig {
	return (*APIConfig)(c)
}

func (c *PortalConfig) API() *APIConfig {
	return (*APIConfig)(c)
}

func (c *SQLConfig) DsnConfig() *mysql.Config {
	return &mysql.Config{
		User:                 c.User,
		Passwd:               c.Pass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", c.Host, c.Port),
		DBName:               c.Name,
		Collation:            "utf8mb4_general_ci",
		ParseTime:            true,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
}

func (c *SQLConfig) DsnConfigWithoutName() *mysql.Config {
	cfg := c.DsnConfig()
	cfg.DBName = ""
	return cfg
}

func (c *SQLConfig) GormConfig() *gorm.Config {
	var (
		logLevel  = logger.Warn
		ignoreRNF = true
	)

	if c.Verbose {
		logLevel = logger.Info
		ignoreRNF = false
	}

	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: ignoreRNF,
				Colorful:                  true,
			},
		),
		NowFunc: func() time.Time {
			return time.Now().Truncate(time.Microsecond)
		},
		TranslateError: true,
	}
}
