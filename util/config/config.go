package config

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	goflag "flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	config Config
	parsed atomic.Bool
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
	parsed.Store(false)

	pflag.Bool("production", false, "whether production or development")
	pflag.Int("port", 1323, "api port")
	pflag.Bool("only-migrate", false, "only migrate db (not start server)")
	pflag.Bool("insert-mock-data", false, "insert sample mock data(for dev)")

	pflag.String("db-user", "", "db user name")
	pflag.String("db-pass", "", "db password")
	pflag.String("db-host", "127.0.0.1", "db host")
	pflag.String("db-name", "", "db name")
	pflag.Int("db-port", 3306, "db port")
	pflag.Bool("db-verbose", false, "db verbose mode")
	pflag.String("traq-cookie", "", "traq cookie")
	pflag.String("traq-api-endpoint", "", "traq api endpoint")
	pflag.String("knoq-cookie", "", "knoq cookie")
	pflag.String("knoq-api-endpoint", "", "knoq api endpoint")
	pflag.String("portal-cookie", "", "portal cookie")
	pflag.String("portal-api-endpoint", "", "portal api endpoint")
	pflag.StringP("config", "c", "", "config file path")
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func Parse() error {
	pflag.Parse()
	if err := ReadFromFile(); err != nil {
		return fmt.Errorf("read config from file: %w", err)
	}

	return nil
}

func ReadFromFile() error {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("bind pflags: %w", err)
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

	parsed.Store(true)

	return nil
}

func GetConfig() *Config {
	if !parsed.Load() {
		panic("config does not parsed")
	}
	return &config
}

func GetModified(editFunc EditFunc) *Config {
	cloned := config.clone()
	editFunc(cloned)
	return cloned
}

func (c *Config) clone() *Config {
	cloned := *c
	return &cloned
}

func (c *Config) IsDevelopment() bool {
	return !c.IsProduction
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *Config) IsOnlyMigrate() bool {
	return c.OnlyMigrate
}

func (c *Config) InsertMock() bool {
	return c.InsertMockData
}

func (c *Config) SQLConf() *SQLConfig {
	return &c.DB
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

func (c *SQLConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&collation=utf8mb4_general_ci", c.User, c.Pass, c.Host, c.Port, c.Name)
}

func (c *SQLConfig) DsnWithoutName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&collation=utf8mb4_general_ci", c.User, c.Pass, c.Host, c.Port)
}

func (c *SQLConfig) DBName() string {
	return c.Name
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
	}
}
