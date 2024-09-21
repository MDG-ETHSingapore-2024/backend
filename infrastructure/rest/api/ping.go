package api

import (
	"backend/infrastructure/repository/db/postgres"
	"backend/infrastructure/repository/db/operations"
	"net/http"
	"strings"
	"fmt"
	"github.com/labstack/echo/v4"
)

func PingPong() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(200, "pong")
	}
}

func getTableId(c echo.Context) string {
	return c.Param("tableId")
}

func getQueryParams(c echo.Context) map[string][]string {
	return c.QueryParams()
}

func parseColumns(queryParams map[string][]string) []string {
	columnsParam := queryParams["columns"]
	if len(columnsParam) > 0 && columnsParam[0] != "" {
		return strings.Split(columnsParam[0], ",")
	}
	return []string{"*"}
}

func GetData(c echo.Context) error {
	tableId := getTableId(c)
	queryParams := getQueryParams(c)
	columns := parseColumns(queryParams)
	ops := parseQueryParams(queryParams)

	query := postgres.ConstructSelectQuery(tableId, columns, ops)
	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, result.Rows())
}

func InsertData(c echo.Context) error {
	tableId := getTableId(c)
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}

	columns, values := prepareColumnsAndValues(data)

	query := postgres.ConstructInsertQuery(tableId, columns, values)
	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, "Data inserted successfully")
}

func UpdateData(c echo.Context) error {
	tableId := getTableId(c)
	queryParams := getQueryParams(c)
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}

	columns, values := prepareColumnsAndValues(data)
	ops := parseQueryParams(queryParams)

	query := postgres.ConstructUpdateQuery(tableId, columns, values, ops)
	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, "Data updated successfully")
}

func DeleteData(c echo.Context) error {
	tableId := getTableId(c)
	queryParams := getQueryParams(c)
	ops := parseQueryParams(queryParams)

	query := postgres.ConstructDeleteQuery(tableId, ops)
	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	result := repo.ExecuteQuery(query)
	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to execute query")
	}

	return c.JSON(http.StatusOK, "Data deleted successfully")
}

func parseQueryParams(queryParams map[string][]string) []operations.OperationObject {
	ops := make([]operations.OperationObject, 0)

	for key, values := range queryParams {
		if key == "columns" {
			continue
		}

		for _, value := range values {
			parts := strings.Split(value, ":")
			if len(parts) != 2 {
				continue
			}

			opCode := parts[0]
			val := parts[1]
			op := operations.GetOperation(opCode)

			ops = append(ops, operations.OperationObject{
				Operation:  op,
				ColumnName: key,
				Value:      fmt.Sprintf("'%s'", val),
			})
		}
	}

	return ops
}

func prepareColumnsAndValues(data map[string]string) ([]string, []string) {
	columns := make([]string, 0, len(data))
	values := make([]string, 0, len(data))

	for k, v := range data {
		columns = append(columns, k)
		values = append(values, fmt.Sprintf("'%s'", v))
	}

	return columns, values
}
