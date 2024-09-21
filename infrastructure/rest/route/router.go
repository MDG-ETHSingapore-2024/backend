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
		router.GET("/:projectId/:tableId", api.GetData)
		router.POST("/:projectId/:tableId", api.InsertData)
		router.PATCH("/:projectId/:tableId", api.UpdateData)
		router.DELETE("/:projectId/:tableId", api.DeleteData)
		router.GET("/ping", api.PingPong())
	}
}
