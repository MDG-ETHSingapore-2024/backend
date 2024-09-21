package postgres

import (
	"backend/infrastructure/repository/db/operations"
	"fmt"
	"log"
	"strings"
)

func ConstructConditions(ops []operations.OperationObject) string {
	stmts := make([]string, len(ops))
    for i, op := range ops { 
        stmts[i] = operations.GetOperationStatement(op)
    }
    return strings.Join(stmts, " AND ")	
}

func ConstructSelectQuery(table string, columns []string, ops []operations.OperationObject) string {
	stmt := ""
	if len(columns) > 0 {
		stmt += fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
	} else {
		stmt += "SELECT * "
	}
	stmt += fmt.Sprintf("FROM %s", table)
	if len(ops) > 0 {
		stmt += fmt.Sprintf("WHERE %s", ConstructConditions(ops))
	}
	stmt += ";"
	return stmt
}

func ConstructInsertQuery(table string, columns []string, values []string) string {
	if len(columns) != len(values) {
		log.Fatalln(fmt.Errorf("Key value lengths does not match"))
		return ""
	}
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, strings.Join(columns, ", "), strings.Join(values, ", "))
	return stmt
}

func ConstructUpdateQuery(table string, columns []string, values []string, ops []operations.OperationObject) string {
	if len(columns) != len(values) {
		log.Fatalln(fmt.Errorf("Key value lengths does not match"))
		return ""
	}
	updates := ""
	for i := range columns {
		updates += fmt.Sprintf("%s = %s", columns[i], values[i])
	}
	stmt := fmt.Sprintf("UPDATE %s SET %s", table, updates)
	if len(ops) > 0 {
		stmt += fmt.Sprintf("WHERE %s", ConstructConditions(ops))
	}
	stmt += ";"
	return stmt
}

func ConstructDeleteQuery(table string, ops []operations.OperationObject) string {
	stmt := fmt.Sprintf("DELETE FROM %s ", table)
	if len(ops) > 0 {
		stmt += fmt.Sprintf("WHERE %s", ConstructConditions(ops))
	}
	stmt += ";"
	return stmt
}