package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
	"github.com/traPtitech/traPortfolio/util/mockdata"
)

func main() {
	migrate := flag.Bool("migrate", false, "migration mode or not")
	flag.Parse()
	if *migrate {
		s := config.SQLConf()
		h, err := infrastructure.NewSQLHandler(&s)
		if err != nil {
			log.Fatal(err)
		}

		if config.IsDevelopment() {
			if err := mockdata.InsertSampleDataToDB(h); err != nil {
				log.Fatal(err)
			}
		}

		log.Println("finished")

		return
	}

	isDevelopment := config.IsDevelopment()
	s := config.SQLConf()
	t := config.TraqConf(isDevelopment)
	p := config.PortalConf(isDevelopment)
	k := config.KnoqConf(isDevelopment)

	api, err := infrastructure.InjectAPIServer(&s, &t, &p, &k)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	if err := handler.Setup(e, api); err != nil {
		log.Fatal(err)
	}

	// Start server
	e.Logger.Fatal(e.Start(config.Port()))
}
