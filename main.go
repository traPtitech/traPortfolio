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
	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}

	appConf := config.Load()
	db, err := infrastructure.NewGormDB(appConf.SQLConf())
	if err != nil {
		log.Fatal(err)
	}

	if appConf.OnlyMigrate {
		log.Println("migration finished")
		return
	}

	if appConf.InsertMockData {
		if appConf.IsProduction {
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
	if err := handler.Setup(e, api, handler.WithRequestLogger()); err != nil {
		log.Fatal(err)
	}

	// Start server
	e.Logger.Fatal(e.Start(appConf.Addr()))
}
