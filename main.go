package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
)

func main() {
	migrate := flag.Bool("migrate", false, "migration mode or not")
	flag.Parse()
	if *migrate {
		conf := sqlConf()
		_, err := infrastructure.NewSQLHandler(&conf)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("finished")
	} else {
		isDevelopment := os.Getenv("APP_ENV") == "development"
		s := sqlConf()
		t := traQConf(isDevelopment)
		p := portalConf(isDevelopment)
		k := knoQConf(isDevelopment)
		g := groupConf(isDevelopment)

		api, err := infrastructure.InjectAPIServer(&s, &t, &p, &k, &g)
		if err != nil {
			log.Fatal(err)
		}

		e := echo.New()
		infrastructure.Setup(e, api)

		port := os.Getenv("PORT")
		if port == "" {
			port = ":1323"
		}
		// Start server
		e.Logger.Fatal(e.Start(port))
	}
}

func sqlConf() infrastructure.SQLConfig {
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "mysql"
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 3306
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "portfolio"
	}

	return infrastructure.NewSQLConfig(user, password, host, dbname, port)
}

func traQConf(isDevelopment bool) infrastructure.TraQConfig {
	traQCookie := os.Getenv("TRAQ_COOKIE")
	traQAPIEndpoint := os.Getenv("TRAQ_API_ENDPOINT")

	return infrastructure.NewTraQConfig(traQCookie, traQAPIEndpoint, isDevelopment)
}

func knoQConf(isDevelopment bool) infrastructure.KnoQConfig {
	knoQCookie := os.Getenv("KNOQ_COOKIE")
	knoQAPIEndpoint := os.Getenv("KNOQ_API_ENDPOINT")

	return infrastructure.NewKnoqConfig(knoQCookie, knoQAPIEndpoint, isDevelopment)
}

func portalConf(isDevelopment bool) infrastructure.PortalConfig {
	portalCookie := os.Getenv("PORTAL_COOKIE")
	portalAPIEndpoint := os.Getenv("PORTAL_API_ENDPOINT")
	return infrastructure.NewPortalConfig(portalCookie, portalAPIEndpoint, isDevelopment)
}

func groupConf() infrastructure.GroupConfig {
	traQCookie := os.Getenv("TRAQ_COOKIE")
	traQAPIEndpoint := os.Getenv("TRAQ_API_ENDPOINT")
	return infrastructure.NewgGoupConfig(traQCookie, traQAPIEndpoint)
}
