package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/traPtitech/traPortfolio/infrastructure"
)

const (
	envAppEnv            = "APP_ENV"
	envPort              = "PORT"
	envDBUser            = "DB_USER"
	envDBPass            = "DB_PASSWORD"
	envDBHost            = "DB_HOST"
	envDBName            = "DB_DATABASE"
	envDBPort            = "DB_PORT"
	envTraqCookie        = "TRAQ_COOKIE"
	envTraqAPIEndpoint   = "TRAQ_API_ENDPOINT"
	envKnoqCookie        = "KNOQ_COOKIE"
	envKnoqAPIEndpoint   = "KNOQ_API_ENDPOINT"
	envPortalCookie      = "PORTAL_COOKIE"
	envPortalAPIEndpoint = "PORTAL_API_ENDPOINT"
	defaultPort          = 1323
	defaultDBUser        = "root"
	defaultDBPass        = "password"
	defaultDBHost        = "mysql"
	defaultDBName        = "portfolio"
	defaultDBPort        = 3306
)

func IsDevelopment() bool {
	return os.Getenv(envAppEnv) == "development"
}

func Port() string {
	return fmt.Sprintf(":%d", getNumEnvOrDefault(envPort, defaultPort))
}

func SQLConf() infrastructure.SQLConfig {
	user := getEnvOrDefault(envDBUser, defaultDBUser)
	pass := getEnvOrDefault(envDBPass, defaultDBPass)
	host := getEnvOrDefault(envDBHost, defaultDBHost)
	dbname := getEnvOrDefault(envDBName, defaultDBName)
	port := getNumEnvOrDefault(envDBPort, defaultDBPort)

	return infrastructure.NewSQLConfig(user, pass, host, dbname, port)
}

func TraqConf(isDevelopment bool) infrastructure.TraQConfig {
	traQCookie := os.Getenv(envTraqCookie)
	traQAPIEndpoint := os.Getenv(envTraqAPIEndpoint)

	return infrastructure.NewTraQConfig(traQCookie, traQAPIEndpoint, isDevelopment)
}

func KnoqConf(isDevelopment bool) infrastructure.KnoQConfig {
	knoQCookie := os.Getenv(envKnoqCookie)
	knoQAPIEndpoint := os.Getenv(envKnoqAPIEndpoint)

	return infrastructure.NewKnoqConfig(knoQCookie, knoQAPIEndpoint, isDevelopment)
}

func PortalConf(isDevelopment bool) infrastructure.PortalConfig {
	portalCookie := os.Getenv(envPortalCookie)
	portalAPIEndpoint := os.Getenv(envPortalAPIEndpoint)

	return infrastructure.NewPortalConfig(portalCookie, portalAPIEndpoint, isDevelopment)
}

func getEnvOrDefault(env string, def string) string {
	s := os.Getenv(env)
	if len(s) == 0 {
		return def
	}

	return s
}

func getNumEnvOrDefault(env string, def int) int {
	i, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return def
	}

	return i
}
