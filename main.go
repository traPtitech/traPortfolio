package main

import (
	"flag"
	"log"
	"os"
	"strconv"

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
		s := sqlConf()
		t := traQConf()
		p := portalConf()
		k := knoQConf()
		infrastructure.Init(&s, &t, &p, &k)
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

func traQConf() infrastructure.TraQConfig {
	traQCookie := os.Getenv("TRAQ_COOKIE")
	traQAPIEndpoint := os.Getenv("TRAQ_API_ENDPOINT")

	return infrastructure.NewTraQConfig(traQCookie, traQAPIEndpoint)
}

func knoQConf() infrastructure.KnoQConfig {
	knoQCookie := os.Getenv("KNOQ_COOKIE")
	knoQAPIEndpoint := os.Getenv("KNOQ_API_ENDPOINT")

	return infrastructure.NewKnoqConfig(knoQCookie, knoQAPIEndpoint)
}

func portalConf() infrastructure.PortalConfig {
	portalCookie := os.Getenv("PORTAL_COOKIE")
	portalAPIEndpoint := os.Getenv("PORTAL_API_ENDPOINT")
	return infrastructure.NewPortalConfig(portalCookie, portalAPIEndpoint)
}
