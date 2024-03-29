package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
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
	Config struct {
		IsProduction   bool `mapstructure:"production"`
		Port           int
		OnlyMigrate    bool
		InsertMockData bool

		DB     SQLConfig
		Traq   APIConfig
		Knoq   APIConfig
		Portal APIConfig
	}

	SQLConfig struct {
		User    string
		Pass    string
		Host    string
		Name    string
		Port    int
		Verbose bool
	}

	APIConfig struct {
		Cookie      string
		APIEndpoint string
	}
)

func init() {
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

type LoadOpts struct {
	SkipReadFromFiles bool
}

func Load(opts LoadOpts) (*Config, error) {
	var c Config

	pflag.Parse()
	for _, key := range flagKeys {
		if err := viper.BindPFlag(key.path, pflag.Lookup(key.flag)); err != nil {
			return nil, fmt.Errorf("bind flag %s: %w", key.flag, err)
		}
	}

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("bind flags: %w", err)
	}

	if !opts.SkipReadFromFiles {
		configPath := viper.GetString("config")
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
					return nil, fmt.Errorf("read config from %s: %w", configPath, err)
				}

				log.Printf("config file does not found: %v\n", err)
			} else {
				return nil, fmt.Errorf("read config: %w", err)
			}
		} else {
			log.Printf("successfully loaded config from %s", viper.ConfigFileUsed())
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &c, nil
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
