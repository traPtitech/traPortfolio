package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Config struct {
		IsProduction   bool `mapstructure:"production"`
		Port           int
		OnlyMigrate    bool
		InsertMockData bool

		DB     SQLConfig
		Traq   TraqConfig
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

	TraqConfig struct {
		AccessToken string
	}
)

func init() {
	pflag.Bool("production", false, "whether production or development")
	viper.BindPFlag("production", pflag.Lookup("production"))

	pflag.Int("port", 1323, "api port")
	viper.BindPFlag("port", pflag.Lookup("port"))

	pflag.Bool("only-migrate", false, "only migrate db (not start server)")
	viper.BindPFlag("onlyMigrate", pflag.Lookup("only-migrate"))

	pflag.Bool("insert-mock-data", false, "insert sample mock data(for dev)")
	viper.BindPFlag("insertMockData", pflag.Lookup("insert-mock-data"))

	pflag.String("db-user", "root", "db user name")
	viper.BindPFlag("db.user", pflag.Lookup("db-user"))

	pflag.String("db-pass", "password", "db password")
	viper.BindPFlag("db.pass", pflag.Lookup("db-pass"))

	pflag.String("db-host", "localhost", "db host")
	viper.BindPFlag("db.host", pflag.Lookup("db-host"))

	pflag.String("db-name", "portfolio", "db name")
	viper.BindPFlag("db.name", pflag.Lookup("db-name"))

	pflag.Int("db-port", 3306, "db port")
	viper.BindPFlag("db.port", pflag.Lookup("db-port"))

	pflag.Bool("db-verbose", false, "db verbose mode")
	viper.BindPFlag("db.verbose", pflag.Lookup("db-verbose"))

	pflag.String("traq-access-token", "", "traq access token")
	viper.BindPFlag("traq.accessToken", pflag.Lookup("traq-access-token"))

	pflag.String("knoq-cookie", "", "knoq cookie")
	viper.BindPFlag("knoq.cookie", pflag.Lookup("knoq-cookie"))

	pflag.String("knoq-api-endpoint", "", "knoq api endpoint")
	viper.BindPFlag("knoq.apiEndpoint", pflag.Lookup("knoq-api-endpoint"))

	pflag.String("portal-cookie", "", "portal cookie")
	viper.BindPFlag("portal.cookie", pflag.Lookup("portal-cookie"))

	pflag.String("portal-api-endpoint", "", "portal api endpoint")
	viper.BindPFlag("portal.apiEndpoint", pflag.Lookup("portal-api-endpoint"))

	pflag.StringP("config", "c", "", "config file path")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

type LoadOpts struct {
	SkipReadFromFiles bool
}

func Load(opts LoadOpts) (*Config, error) {
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("bind flags: %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("TPF")
	viper.AutomaticEnv()

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

				log.Printf("config file did not used: %v\n", err)
			} else {
				return nil, fmt.Errorf("read config: %w", err)
			}
		} else {
			log.Printf("successfully loaded config from %s", viper.ConfigFileUsed())
		}
	}

	var c Config
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
