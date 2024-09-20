package route

import (
	"backend/api"

	"github.com/labstack/echo/v4"
)

func InitRouter() {
	e := echo.New()
	e.GET("/ping", api.PingPong())
	e.Logger.Fatal(e.Start(":1323"))
}
