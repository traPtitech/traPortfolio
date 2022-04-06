package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
)

func main() {
	config.Parse()
	appConf := config.GetConfig()

	if appConf.IsMigrate() {
		s := appConf.SQLConf()
		_, err := infrastructure.NewSQLHandler(&s)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("finished")
	} else {
		s := appConf.SQLConf()
		t := appConf.TraqConf()
		p := appConf.PortalConf()
		k := appConf.KnoqConf()

		api, err := infrastructure.InjectAPIServer(&s, &t, &p, &k)
		if err != nil {
			log.Fatal(err)
		}

		e := echo.New()
		if err := handler.Setup(e, api); err != nil {
			log.Fatal(err)
		}

		// Start server
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", appConf.Port)))
	}
}
