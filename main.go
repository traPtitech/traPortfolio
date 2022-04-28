package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
	"github.com/traPtitech/traPortfolio/util/mockdata"
)

func main() {
	config.Parse()
	appConf := config.GetConfig()

	if appConf.IsMigrate() {
		s := appConf.SQLConf()
		h, err := infrastructure.NewSQLHandler(s)
		if err != nil {
			log.Fatal(err)
		}

		if appConf.IsDevelopment() {
			if err := mockdata.InsertSampleDataToDB(h); err != nil {
				log.Fatal(err)
			}
		}

		log.Println("finished")

		return
	}

	api, err := infrastructure.InjectAPIServer(appConf, appConf.IsDevelopment())
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	if err := handler.Setup(e, api); err != nil {
		log.Fatal(err)
	}

	// Start server
	e.Logger.Fatal(e.Start(appConf.Addr()))
}
