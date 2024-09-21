package api

import (
	"fmt"
	"net/http"
	"strings"

	"backend/infrastructure/repository/db/operations"
	"backend/infrastructure/repository/db/postgres"

	"github.com/labstack/echo/v4"
)

func getUserIdFromContext(c echo.Context) string {
	return c.Request().Header.Get("X-Wallet-Address")
}

func CreateDatabase(c echo.Context) error {
	userId := getUserIdFromContext(c)

	type RequestBody struct {
		DbName string `json:"dbName"`
		Type   string `json:"type"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if userId == "" || body.DbName == "" || body.Type == "" {
		return c.JSON(http.StatusBadRequest, "Missing userId, dbName, or type")
	}

	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	createRoleQuery := fmt.Sprintf("DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s') THEN CREATE ROLE \"%s\" LOGIN; END IF; END $$;", userId, userId)
	roleResult := repo.ExecuteQuery(createRoleQuery)
	if roleResult.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create role for userId: %s", userId))
	}

	var createDbQuery string
	switch body.Type {
	case "sql":
		createDbQuery = fmt.Sprintf("CREATE DATABASE \"%s\" OWNER \"%s\";", body.DbName, userId)
	case "nosql":
		createDbQuery = fmt.Sprintf("CREATE DATABASE \"%s\" OWNER \"%s\";", body.DbName, userId)
	default:
		return c.JSON(http.StatusBadRequest, "Invalid database type. Must be 'sql' or 'nosql'.")
	}

	dbResult := repo.ExecuteQuery(createDbQuery)
	if dbResult.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create database: %s", body.DbName))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Database %s of type %s created successfully with owner %s", body.DbName, body.Type, userId))
}

func RetrieveAllDatabases(c echo.Context) error {
	userId := getUserIdFromContext(c)

	if userId == "" {
		return c.JSON(http.StatusBadRequest, "Missing userId")
	}

	repo := postgres.OpenDatabase("postgres", operations.ConnectionString)
	defer repo.CloseDatabase()

	query := fmt.Sprintf(`
        SELECT datname 
        FROM pg_database 
        WHERE datdba = (SELECT oid FROM pg_roles WHERE rolname = '%s') 
        AND datistemplate = false;`, userId)

	result := repo.ExecuteQuery(query)

	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve databases")
	}

	dbNames := []string{}
	rows := result.Rows()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to scan database name")
		}
		dbNames = append(dbNames, dbName)
	}

	return c.JSON(http.StatusOK, dbNames)
}

func CreateTable(c echo.Context) error {
	dbName := c.Param("dbName")
	userId := getUserIdFromContext(c)

	type RequestBody struct {
		TableName   string   `json:"tableName"`
		ColumnNames []string `json:"columnNames"`
		ColumnTypes []string `json:"columnTypes"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if dbName == "" || body.TableName == "" || len(body.ColumnNames) == 0 || len(body.ColumnNames) != len(body.ColumnTypes) {
		return c.JSON(http.StatusBadRequest, "Missing dbName, tableName, or columns, or columnNames and columnTypes length mismatch")
	}

	connStr := fmt.Sprintf("%s user=%s dbname=%s", operations.ConnectionString, userId, dbName)
	repo := postgres.OpenDatabase("postgres", connStr)
	defer repo.CloseDatabase()

	columnDefs := []string{}
	for i := range body.ColumnNames {
		columnDefs = append(columnDefs, fmt.Sprintf("\"%s\" %s", body.ColumnNames[i], body.ColumnTypes[i]))
	}

	query := fmt.Sprintf("CREATE TABLE \"%s\" (%s);", body.TableName, strings.Join(columnDefs, ", "))
	result := repo.ExecuteQuery(query)

	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to create table")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Table %s created successfully in %s with columns %v", body.TableName, dbName, body.ColumnNames))
}


func AddColumn(c echo.Context) error {
	dbName := c.Param("dbName")
	tableName := c.Param("tableName")
	userId := getUserIdFromContext(c)

	type RequestBody struct {
		ColumnName string `json:"columnName"`
		ColumnType string `json:"columnType"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if dbName == "" || tableName == "" || body.ColumnName == "" || body.ColumnType == "" {
		return c.JSON(http.StatusBadRequest, "Missing dbName, tableName, columnName, or columnType")
	}

	connStr := fmt.Sprintf("%s user=%s dbname=%s", operations.ConnectionString, userId, dbName)
	repo := postgres.OpenDatabase("postgres", connStr)
	defer repo.CloseDatabase()

	query := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN \"%s\" %s;", tableName, body.ColumnName, body.ColumnType)
	result := repo.ExecuteQuery(query)

	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to add column")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Column %s added successfully to table %s", body.ColumnName, tableName))
}

func AddRow(c echo.Context) error {
	dbName := c.Param("dbName")
	tableName := c.Param("tableName")
	userId := getUserIdFromContext(c)

	var rowData map[string]interface{}
	if err := c.Bind(&rowData); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	delete(rowData, "dbName")
	delete(rowData, "tableName")

	if dbName == "" || tableName == "" || len(rowData) == 0 {
		return c.JSON(http.StatusBadRequest, "Missing dbName, tableName, or rowData")
	}

	connStr := fmt.Sprintf("%s user=%s dbname=%s", operations.ConnectionString, userId, dbName)
	repo := postgres.OpenDatabase("postgres", connStr)
	defer repo.CloseDatabase()

	columns := []string{}
	values := []string{}
	for col, val := range rowData {
		columns = append(columns, fmt.Sprintf("\"%s\"", col))
		values = append(values, fmt.Sprintf("'%v'", val))
	}

	query := fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s);", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
	result := repo.ExecuteQuery(query)

	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to add row")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Row added successfully to table %s in database %s", tableName, dbName))
}

func ViewAllTables(c echo.Context) error {
	dbName := c.Param("dbName")
	userId := getUserIdFromContext(c)

	if dbName == "" {
		return c.JSON(http.StatusBadRequest, "Missing dbName")
	}

	connStr := fmt.Sprintf("%s user=%s dbname=%s", operations.ConnectionString, userId, dbName)
	repo := postgres.OpenDatabase("postgres", connStr)
	defer repo.CloseDatabase()

	query := `
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public'
        ORDER BY table_name;`

	result := repo.ExecuteQuery(query)

	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve tables")
	}

	tables := []string{}
	rows := result.Rows()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to scan table name")
		}
		tables = append(tables, tableName)
	}

	return c.JSON(http.StatusOK, tables)
}

func ViewAllData(c echo.Context) error {
	dbName := c.Param("dbName")
	tableName := c.Param("tableName")
	userId := getUserIdFromContext(c)

	if dbName == "" || tableName == "" {
		return c.JSON(http.StatusBadRequest, "Missing dbName or tableName")
	}

	connStr := fmt.Sprintf("%s user=%s dbname=%s", operations.ConnectionString, userId, dbName)
	repo := postgres.OpenDatabase("postgres", connStr)
	defer repo.CloseDatabase()

	query := fmt.Sprintf("SELECT * FROM \"%s\";", tableName)
	result := repo.ExecuteQuery(query)

	if result.Status() == postgres.Failure {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve data")
	}

	columns, _ := result.Rows().Columns()
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	data := []map[string]interface{}{}

	for result.Rows().Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		result.Rows().Scan(valuePtrs...)

		rowData := map[string]interface{}{}
		for i, col := range columns {
			rowData[col] = values[i]
		}

		data = append(data, rowData)
	}

	return c.JSON(http.StatusOK, data)
}
