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
		router.GET("/dbs", api.RetrieveAllDatabases)
		router.GET("/:dbName/tables", api.ViewAllTables)
		router.GET("/:dbName/:tableName/data", api.ViewAllData)
		router.POST("/db", api.CreateDatabase)
		router.POST("/:dbName/table", api.CreateTable)
		router.POST("/:dbName/:tableName/column", api.AddColumn)
		router.POST("/:dbName/:tableName/row", api.AddRow)

		// NoSQL Collection Operations
		// router.POST("/:dbId/collection", api.CreateTable)
		// router.POST("/:dbId/:collectionId/document", api.AddDocument)
		// router.POST("/:dbId/:collectionId/attribute", api.AddAttribute)
	}
}
