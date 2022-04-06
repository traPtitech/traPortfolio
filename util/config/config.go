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
	defaultAppEnv  = "development"
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

// ReadOnly outside this package
type Config struct {
	AppEnv            string `mapstructure:"appEnv"`
	Port              int    `mapstructure:"port"`
	DBUser            string `mapstructure:"dbUser"`
	DBPass            string `mapstructure:"dbPass"`
	DBHost            string `mapstructure:"dbHost"`
	DBName            string `mapstructure:"dbName"`
	DBPort            int    `mapstructure:"dbPort"`
	TraqCookie        string `mapstructure:"traqCookie"` // r_session
	TraqAPIEndpoint   string `mapstructure:"traqAPIEndpoint"`
	KnoqCookie        string `mapstructure:"knoqCookie"` // session
	KnoqAPIEndpoint   string `mapstructure:"knoqAPIEndpoint"`
	PortalCookie      string `mapstructure:"portalCookie"`
	PortalAPIEndpoint string `mapstructure:"portalAPIEndpoint"`
	Migrate           bool   `mapstructure:"migrate"`
}

func init() {
	pflag.String("app-env", defaultAppEnv, "whether production of development")
	pflag.Int("port", defaultAppPort, "api port")
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
	pflag.String("config-dir-path", "", "config directory path")
	pflag.Bool("migrate", false, "run with migrate mode")
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func Parse() {
	pflag.Parse()
	setParsed(true)

	_ = viper.BindPFlag("appEnv", pflag.Lookup("app-env"))
	_ = viper.BindPFlag("port", pflag.Lookup("port"))
	_ = viper.BindPFlag("dbUser", pflag.Lookup("db-user"))
	_ = viper.BindPFlag("dbPass", pflag.Lookup("db-pass"))
	_ = viper.BindPFlag("dbHost", pflag.Lookup("db-host"))
	_ = viper.BindPFlag("dbName", pflag.Lookup("db-name"))
	_ = viper.BindPFlag("dbPort", pflag.Lookup("db-port"))
	_ = viper.BindPFlag("traqCookie", pflag.Lookup("traq-cookie"))
	_ = viper.BindPFlag("traqAPIEndpoint", pflag.Lookup("traq-api-endpoint"))
	_ = viper.BindPFlag("knoqCookie", pflag.Lookup("knoq-cookie"))
	_ = viper.BindPFlag("knoqAPIEndpoint", pflag.Lookup("knoq-api-endpoint"))
	_ = viper.BindPFlag("portalCookie", pflag.Lookup("portal-cookie"))
	_ = viper.BindPFlag("portalAPIEndpoint", pflag.Lookup("portal-api-endpoint"))

	_ = viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")

	configPath, err := pflag.CommandLine.GetString("config-dir-path")
	if err != nil {
		panic(err)
	}
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Printf("config file does not found %s", err.Error())
		} else {
			// Config file was found but another error was produced
			log.Fatal("cannot load config", err)
		}
	} else {
		log.Printf("successfully loaded config")
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	if !isParsed() {
		panic("config does not parsed")
	}
	return &config
}

func (c *Config) SetDBName(name string) *Config {
	cloned := c.Clone()
	cloned.DBName = name
	return &cloned
}

func (c *Config) Clone() Config {
	return *c
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) SQLConf() infrastructure.SQLConfig {
	return infrastructure.NewSQLConfig(c.DBUser, c.DBPass, c.DBHost, c.DBName, c.DBPort)
}

func (c *Config) TraqConf() infrastructure.TraQConfig {
	return infrastructure.NewTraQConfig(c.TraqCookie, c.TraqAPIEndpoint, c.IsDevelopment())
}

func (c *Config) KnoqConf() infrastructure.KnoQConfig {
	return infrastructure.NewKnoqConfig(c.KnoqCookie, c.KnoqAPIEndpoint, c.IsDevelopment())
}

func (c *Config) PortalConf() infrastructure.PortalConfig {
	return infrastructure.NewPortalConfig(c.PortalCookie, c.PortalAPIEndpoint, c.IsDevelopment())
}

func (c *Config) IsMigrate() bool { return c.Migrate }
