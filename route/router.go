package route

import (
	"backend/api"
	"backend/middleware"

	"github.com/labstack/echo/v4"
)

func InitRouter() {
	e := echo.New()
	e.Use(middleware.SessionAuth)
	e.GET("/ping", api.PingPong())
	e.Logger.Fatal(e.Start(":1323"))
}
