package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/internal/handler"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/internal/pkgs/config"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
)

func main() {
	appConf, err := config.Load(config.LoadOpts{})
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewGormDB(appConf.DB)
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

		if err := mockdata.InsertSampleDataToDB(db); err != nil {
			log.Fatal(err)
		}
	}

	api, err := injectIntoAPIServer(appConf, db)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	if err := handler.Setup(appConf.IsProduction, e, api, handler.WithRequestLogger()); err != nil {
		log.Fatal(err)
	}

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", appConf.Port)))
}
