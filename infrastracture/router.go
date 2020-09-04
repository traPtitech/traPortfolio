package infrastructure

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/traPtitech/traPortfolio/di"
)

func Init() {
	// Echo instance
	e := echo.New()

	// userController := controllers.NewUserController(NewSqlHandler())
	// pingRepository := repository.NewPingRepository()
	// pingInteractor := interactor.NewPingInteractor(pingRepository)
	// pingController := controllers.NewPingController(pingInteractor)
	api := di.InjectAPIServer()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.GET("/users", func(c echo.Context) error { return userController.Index(c) })
	// e.GET("/user/:id", func(c echo.Context) error { return userController.Show(c) })
	// e.POST("/create", func(c echo.Context) error { return userController.Create(c) })
	// e.PUT("/user/:id", func(c echo.Context) error { return userController.Save(c) })
	// e.DELETE("/user/:id", func(c echo.Context) error { return userController.Delete(c) })
	e.GET("/ping", func(c echo.Context) error { return api.Ping.Ping(c) })

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
