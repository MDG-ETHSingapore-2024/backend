package route

import (
	"backend/infrastructure/rest/api"
	"backend/infrastructure/rest/middleware"

	"github.com/labstack/echo/v4"
)

func InitRouter(e *echo.Echo) {
	router := e.Group("/v1")
	router.Use(middleware.SessionAuth)
	{
		router.GET("/ping", api.PingPong())
	}
}
