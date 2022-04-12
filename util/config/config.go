package config

import (
	"log"
	"sync"

	goflag "flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/traPtitech/traPortfolio/infrastructure"
)

const (
	defaultAppPort = 1323
	defaultDBPort  = 3306
	defaultDBHost  = "127.0.0.1"
)

var (
	config Config
	parsed bool
	rmu    sync.RWMutex
)

func setParsed(b bool) {
	rmu.Lock()
	defer rmu.Unlock()
	parsed = b
}

func isParsed() bool {
	rmu.RLock()
	defer rmu.RUnlock()
	return parsed
}

// Immutable except within this package or EditFunc
type Config struct {
	IsProduction bool `mapstructure:"production"`
	Port         int  `mapstructure:"port"`
	Migrate      bool `mapstructure:"migrate"`

	DB struct {
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
		Host string `mapstructure:"host"`
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"db"`

	Traq struct {
		Cookie      string `mapstructure:"cookie"` // r_session
		APIEndpoint string `mapstructure:"apiEndpoint"`
	} `mapstructure:"traq"`

	Knoq struct {
		Cookie      string `mapstructure:"cookie"` // session
		APIEndpoint string `mapstructure:"apiEndpoint"`
	} `mapstructure:"knoq"`

	Portal struct {
		Cookie      string `mapstructure:"cookie"` // access_token
		APIEndpoint string `mapstructure:"apiEndpoint"`
	} `mapstructure:"portal"`
}

func init() {
	pflag.Bool("production", false, "whether production or development")
	pflag.Int("port", defaultAppPort, "api port")
	pflag.Bool("migrate", false, "run with migrate mode")

	pflag.String("db-user", "", "db user name")
	pflag.String("db-pass", "", "db password")
	pflag.String("db-host", defaultDBHost, "db host")
	pflag.String("db-name", "", "db name")
	pflag.Int("db-port", defaultDBPort, "db port")
	pflag.String("traq-cookie", "", "traq cookie")
	pflag.String("traq-api-endpoint", "", "traq api endpoint")
	pflag.String("knoq-cookie", "", "knoq cookie")
	pflag.String("knoq-api-endpoint", "", "knoq api endpoint")
	pflag.String("portal-cookie", "", "portal cookie")
	pflag.String("portal-api-endpoint", "", "portal api endpoint")
	pflag.StringP("config", "c", "", "config file path")
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func Parse() {
	pflag.Parse()

	_ = viper.BindPFlag("production", pflag.Lookup("isProduction"))
	_ = viper.BindPFlag("port", pflag.Lookup("port"))
	_ = viper.BindPFlag("migration", pflag.Lookup("migration"))

	_ = viper.BindPFlag("db.user", pflag.Lookup("db-user"))
	_ = viper.BindPFlag("db.pass", pflag.Lookup("db-pass"))
	_ = viper.BindPFlag("db.host", pflag.Lookup("db-host"))
	_ = viper.BindPFlag("db.name", pflag.Lookup("db-name"))
	_ = viper.BindPFlag("db.port", pflag.Lookup("db-port"))
	_ = viper.BindPFlag("traq.cookie", pflag.Lookup("traq-cookie"))
	_ = viper.BindPFlag("traq.apiEndpoint", pflag.Lookup("traq-api-endpoint"))
	_ = viper.BindPFlag("knoq.cookie", pflag.Lookup("knoq-cookie"))
	_ = viper.BindPFlag("knoq.apiEndpoint", pflag.Lookup("knoq-api-endpoint"))
	_ = viper.BindPFlag("portal.cookie", pflag.Lookup("portal-cookie"))
	_ = viper.BindPFlag("portal.apiEndpoint", pflag.Lookup("portal-api-endpoint"))

	_ = viper.BindPFlags(pflag.CommandLine)

	configPath, err := pflag.CommandLine.GetString("config")
	if err != nil {
		log.Fatal(err)
	}
	if len(configPath) > 0 {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config") // name of config file (without extension)
		viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired

			// exit if configPath is explicitly specified and fails to load.
			if len(configPath) > 0 {
				log.Fatal("cannot load config from ", configPath)
			}

			log.Printf("config file does not found %s", err.Error())
		} else {
			// Config file was found but another error was produced
			log.Fatal("cannot load config", err)
		}
	} else {
		log.Printf("successfully loaded config from %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	setParsed(true)
}

func GetConfig() *Config {
	if !isParsed() {
		panic("config does not parsed")
	}
	return &config
}

type EditFunc func(*Config)

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

func (c *Config) SQLConf() infrastructure.SQLConfig {
	return infrastructure.NewSQLConfig(c.DB.User, c.DB.Pass, c.DB.Host, c.DB.Name, c.DB.Port)
}

func (c *Config) TraqConf() infrastructure.TraQConfig {
	return infrastructure.NewTraQConfig(c.Traq.Cookie, c.Traq.APIEndpoint, c.IsDevelopment())
}

func (c *Config) KnoqConf() infrastructure.KnoQConfig {
	return infrastructure.NewKnoqConfig(c.Knoq.Cookie, c.Knoq.APIEndpoint, c.IsDevelopment())
}

func (c *Config) PortalConf() infrastructure.PortalConfig {
	return infrastructure.NewPortalConfig(c.Portal.Cookie, c.Portal.APIEndpoint, c.IsDevelopment())
}

func (c *Config) IsMigrate() bool { return c.Migrate }
