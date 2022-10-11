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
	db, err := infrastructure.NewGormDB(appConf.SQLConf())
	if err != nil {
		log.Fatal(err)
	}

	if appConf.IsMigrate() {
		log.Println("migration finished")
		return
	}

	if appConf.InsertMock() {
		if !appConf.IsDevelopment() {
			log.Fatal("cannot specify both `production` and `insert-mock-data`")
		}

		if err := mockdata.InsertSampleDataToDB(infrastructure.NewSQLHandler(db)); err != nil {
			log.Fatal(err)
		}
	}

	api, err := infrastructure.InjectAPIServer(appConf, db)
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
