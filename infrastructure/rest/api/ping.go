package api

import (
	"backend/infrastructure/repository/db/postgres"
	"backend/infrastructure/repository/db/operations"
	"net/http"
	"strings"
	
	"github.com/labstack/echo/v4"
)

func PingPong() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(200, "pong")
	}
}

func GetData(c echo.Context) error {
	projectId := c.Param("projectId")
	tableId := c.Param("tableId")
	queryParams := c.QueryParams()

	columns := []string{"*"}
	ops := parseQueryParams(queryParams)
	query := postgres.ConstructSelectQuery(tableId, columns, ops)

	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.status == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, result.rows)
}

func InsertData(c echo.Context) error {
	projectId := c.Param("projectId")
	tableId := c.Param("tableId")
	var data map[string]string
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}

	columns := make([]string, 0, len(data))
	values := make([]string, 0, len(data))
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
	}

	query := postgres.ConstructInsertQuery(tableId, columns, values)

	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.status == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, "Data inserted successfully")
}

func UpdateData(c echo.Context) error {
	projectId := c.Param("projectId")
	tableId := c.Param("tableId")
	queryParams := c.QueryParams()
	var data map[string]string
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}

	columns := make([]string, 0, len(data))
	values := make([]string, 0, len(data))
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
	}

	ops := parseQueryParams(queryParams)
	query := postgres.ConstructUpdateQuery(tableId, columns, values, ops)

	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.status == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, "Data updated successfully")
}

func DeleteData(c echo.Context) error {
	projectId := c.Param("projectId")
	tableId := c.Param("tableId")
	queryParams := c.QueryParams()

	ops := parseQueryParams(queryParams)
	query := postgres.ConstructDeleteQuery(tableId, ops)

	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.status == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, "Data deleted successfully")
}

func parseQueryParams(queryParams map[string][]string) []operations.OperationObject {
	ops := make([]operations.OperationObject, 0)

	for key, values := range queryParams {
		for _, value := range values {
			// Split the value by ':' to get the operation and the actual value
			parts := strings.Split(value, ":")
			if len(parts) != 2 {
				continue // Skip if the format is incorrect
			}

			opCode := parts[0] // eq, gt, lt, etc.
			val := parts[1]    // actual value

			// Get the operation type from the map
			op := operations.GetOperation(opCode)

			ops = append(ops, operations.OperationObject{
				Operation:  op,
				ColumnName: key,
				Value:      val,
			})
		}
	}

	return ops
}
